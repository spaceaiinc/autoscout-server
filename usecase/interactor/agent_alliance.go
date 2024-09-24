package interactor

import (
	"errors"
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentAllianceInteractor interface {
	// 汎用系 API
	RequestAlliance(input RequestAllianceInput) (RequestAllianceOutput, error)
	UpdateOrCreateAgentAllianceList(input UpdateOrCreateAgentAllianceListInput) (UpdateOrCreateAgentAllianceListOutput, error)
	UpdateAgentAllianceCancelRequest(input UpdateAgentAllianceCancelRequestInput) (UpdateAgentAllianceCancelRequestOutput, error)
	GetAgentAllianceByID(input GetAgentAllianceByIDInput) (GetAgentAllianceByIDOutput, error)
	GetAgentAllianceListByAgentID(input GetAgentAllianceListByAgentIDInput) (GetAgentAllianceListByAgentIDOutput, error)
	CheckAnyAllianceWithoutApplication(input CheckAnyAllianceWithoutApplicationInput) (CheckAnyAllianceWithoutApplicationOutput, error)
}

type AgentAllianceInteractorImpl struct {
	firebase                     usecase.Firebase
	sendgrid                     config.Sendgrid
	agentAllianceRepository      usecase.AgentAllianceRepository
	agentRepository              usecase.AgentRepository
	chatGroupWithAgentRepository usecase.ChatGroupWithAgentRepository
	taskRepository               usecase.TaskRepository
}

// AgentAllianceInteractorImpl is an implementation of AgentAllianceInteractor
func NewAgentAllianceInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	aaR usecase.AgentAllianceRepository,
	caR usecase.AgentRepository,
	cgaR usecase.ChatGroupWithAgentRepository,
	tR usecase.TaskRepository,
) AgentAllianceInteractor {
	return &AgentAllianceInteractorImpl{
		firebase:                     fb,
		sendgrid:                     sg,
		agentAllianceRepository:      aaR,
		agentRepository:              caR,
		chatGroupWithAgentRepository: cgaR,
		taskRepository:               tR,
	}
}

/****************************************************************************************/
// 汎用系 API
//

type RequestAllianceInput struct {
	Param entity.CreateAgentAllianceParam
}

type RequestAllianceOutput struct {
	Alliance *entity.AgentAlliance
}

func (i *AgentAllianceInteractorImpl) RequestAlliance(input RequestAllianceInput) (RequestAllianceOutput, error) {
	var (
		output           RequestAllianceOutput
		newAgentAlliance *entity.AgentAlliance
		agent1ID         uint // 数字が小さい方
		agent2ID         uint // 数字が大きい方
		agent1Request    bool // 数字が小さい方
		agent2Request    bool // 数字が大きい方
	)

	/********* 自社の情報を取得 **********/

	requestAgent, err := i.agentRepository.FindByUUID(input.Param.RequestUUID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			return output, wrapped
		} else {
			fmt.Println(err)
			return output, err
		}
	}

	/********* 自社もしくはたエージェントがすでに申請済みかを確認 **********/

	// agent1_idに数字が小さい方のidを入れるための処理
	if input.Param.MyAgentID < requestAgent.ID {
		agent1ID = input.Param.MyAgentID
		agent2ID = requestAgent.ID
		agent1Request = input.Param.Agent1Request
		agent2Request = input.Param.Agent2Request
	} else {
		agent1ID = requestAgent.ID
		agent2ID = input.Param.MyAgentID
		agent1Request = input.Param.Agent2Request
		agent2Request = input.Param.Agent1Request
	}

	agentAlliance, err := i.agentAllianceRepository.FindByAgentID(agent1ID, agent2ID) // 先に申請されているかを確認
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// まだレコードが作成されていない場合
			newAgentAlliance = entity.NewAgentAlliance(
				agent1ID,
				agent2ID,
				agent1Request,
				agent2Request,
				false, // agent1CancelRequest,
				false, // agent2CancelRequest,
			)
			err := i.agentAllianceRepository.Create(newAgentAlliance)
			if err != nil {
				return output, err
			}
		} else {
			fmt.Println(err)
			return output, err
		}
	} else {
		// すでにレコードが存在する場合

		// 申請中の場合
		if (agentAlliance.Agent1ID == input.Param.MyAgentID && agentAlliance.Agent1Request && !agentAlliance.Agent2Request) ||
			(agentAlliance.Agent2ID == input.Param.MyAgentID && agentAlliance.Agent2Request && !agentAlliance.Agent1Request) {
			return output, errors.New("現在申請中ののエージェントです")
		}

		// アライアンス締結済みの場合
		if agentAlliance.Agent1Request && agentAlliance.Agent2Request {
			return output, errors.New("既にアライアンス締結済みのエージェントです")
		}

		if input.Param.MyAgentID == agentAlliance.Agent1ID {
			// 未申請の場合
			newAgentAlliance = entity.NewAgentAlliance(
				agentAlliance.Agent1ID,
				agentAlliance.Agent2ID,
				true,
				agentAlliance.Agent2Request,
				false, // agent1CancelRequest,
				false, // agent2CancelRequest,
			)
		} else {
			newAgentAlliance = entity.NewAgentAlliance(
				agentAlliance.Agent1ID,
				agentAlliance.Agent2ID,
				agentAlliance.Agent1Request,
				true,
				false, // agent1CancelRequest,
				false, // agent2CancelRequest,
			)
		}

		err := i.agentAllianceRepository.Update(agentAlliance.ID, newAgentAlliance)
		if err != nil {
			return output, err
		}
		newAgentAlliance.ID = agentAlliance.ID
	}

	/********* どちらも申請した場合はチャット作成 **********/

	if newAgentAlliance.Agent1Request && newAgentAlliance.Agent2Request {
		fmt.Println("エージェントのチャットグループを作成")

		_, err := i.chatGroupWithAgentRepository.FindByAgentID(newAgentAlliance.Agent1ID, newAgentAlliance.Agent2ID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				// チャットグループ未作成の場合はグループ作成
				chatGroup := entity.NewChatGroupWithAgent(newAgentAlliance.Agent1ID, newAgentAlliance.Agent2ID)

				err = i.chatGroupWithAgentRepository.Create(chatGroup)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.Alliance = newAgentAlliance

	return output, nil
}

type UpdateOrCreateAgentAllianceListInput struct {
	Param entity.MyAgentIDAndOtherAgentIDListParam
}

type UpdateOrCreateAgentAllianceListOutput struct {
	AgentAllianceList []*entity.AgentAlliance
}

func (i *AgentAllianceInteractorImpl) UpdateOrCreateAgentAllianceList(input UpdateOrCreateAgentAllianceListInput) (UpdateOrCreateAgentAllianceListOutput, error) {
	var (
		output       UpdateOrCreateAgentAllianceListOutput
		allianceList []*entity.AgentAlliance
	)

	for _, otherAgentID := range input.Param.OtherAgentIDList {
		var (
			agent1ID      uint // 数字が小さい方
			agent2ID      uint // 数字が大きい方
			agent1Request bool // 数字が小さい方
			agent2Request bool // 数字が大きい方
		)

		// agent1_idに数字が小さい方のidを入れるための処理
		if input.Param.MyAgentID < otherAgentID {
			agent1ID = input.Param.MyAgentID
			agent2ID = otherAgentID
			agent1Request = true
			agent2Request = false
		} else {
			agent1ID = otherAgentID
			agent2ID = input.Param.MyAgentID
			agent1Request = false
			agent2Request = true
		}

		agentAlliance, err := i.agentAllianceRepository.FindByAgentID(input.Param.MyAgentID, otherAgentID) // 先に申請されているかを確認
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				// まだレコードが作成されていない場合
				// 自分だけ申請完了した状態でレコードを登録
				newAgentAlliance := entity.NewAgentAlliance(
					agent1ID,
					agent2ID,
					agent1Request,
					agent2Request,
					false, // agent1CancelRequest,
					false, // agent2CancelRequest,
				)

				err := i.agentAllianceRepository.Create(newAgentAlliance)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				allianceList = append(allianceList, newAgentAlliance)
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			// すでにレコードがある場合
			if agentAlliance.Agent1Request && agentAlliance.Agent2Request {
				// すでに締結ずみの場合はTrueを返して終了
				allianceList = append(allianceList, agentAlliance)
			} else {
				// 未締結の場合は自身の申請状況を更新
				err := i.agentAllianceRepository.UpdateAgentRequest(agentAlliance.ID, input.Param.MyAgentID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				// レスポンスで正しい値を返すためんいagentAllianceの値を更新
				if input.Param.MyAgentID == agentAlliance.Agent1ID {
					agentAlliance.Agent1Request = true
				} else {
					agentAlliance.Agent2Request = true
				}

				allianceList = append(allianceList, agentAlliance)
			}
		}
	}

	output.AgentAllianceList = allianceList

	return output, nil
}

type UpdateAgentAllianceCancelRequestInput struct {
	Param entity.UpdateAgentAllianceCancelRequestParam
}

type UpdateAgentAllianceCancelRequestOutput struct {
	AgentAlliance     *entity.AgentAlliance
	RemainingTaskList []*entity.Task
}

func (i *AgentAllianceInteractorImpl) UpdateAgentAllianceCancelRequest(input UpdateAgentAllianceCancelRequestInput) (UpdateAgentAllianceCancelRequestOutput, error) {
	var (
		output UpdateAgentAllianceCancelRequestOutput
	)

	agentAlliance, err := i.agentAllianceRepository.FindByID(input.Param.AgentAllianceID)
	if err != nil {
		return output, err
	}

	// タスク残りを確認
	taskList, err := i.taskRepository.GetByEachAgentID(agentAlliance.Agent1ID, agentAlliance.Agent2ID)
	if err != nil {
		return output, err
	}

	if len(taskList) > 0 {
		fmt.Println("タスクが残っているため、アライアンスを解除できません", taskList)
		for _, task := range taskList {
			fmt.Println(task.ID)
		}
		output.RemainingTaskList = taskList
		return output, nil
	}

	// 申請をキャンセルしたエージェントの解除申請をtrueに
	if input.Param.MyAgentID == agentAlliance.Agent1ID {
		agentAlliance.Agent1CancelRequest = true
	} else if input.Param.MyAgentID == agentAlliance.Agent2ID {
		agentAlliance.Agent2CancelRequest = true
	}

	if agentAlliance.Agent1CancelRequest && agentAlliance.Agent2CancelRequest {
		// 両方が解除申請がtrueの場合は、アライアンス申請をfalseに
		agentAlliance.Agent1Request = false
		agentAlliance.Agent2Request = false
	}

	err = i.agentAllianceRepository.Update(input.Param.AgentAllianceID, agentAlliance)
	if err != nil {
		return output, err
	}

	output.AgentAlliance = agentAlliance

	return output, nil
}

type GetAgentAllianceByIDInput struct {
	AgentAllianceID uint
}

type GetAgentAllianceByIDOutput struct {
	AgentAlliance *entity.AgentAlliance
}

func (i *AgentAllianceInteractorImpl) GetAgentAllianceByID(input GetAgentAllianceByIDInput) (GetAgentAllianceByIDOutput, error) {
	var (
		output GetAgentAllianceByIDOutput
	)

	agentAlliance, err := i.agentAllianceRepository.FindByID(input.AgentAllianceID)
	if err != nil {
		return output, err
	}

	output.AgentAlliance = agentAlliance

	return output, nil
}

type GetAgentAllianceListByAgentIDInput struct {
	AgentID uint
}

type GetAgentAllianceListByAgentIDOutput struct {
	AgentAllianceList []*entity.AgentAlliance
}

func (i *AgentAllianceInteractorImpl) GetAgentAllianceListByAgentID(input GetAgentAllianceListByAgentIDInput) (GetAgentAllianceListByAgentIDOutput, error) {
	var (
		output GetAgentAllianceListByAgentIDOutput
	)

	agentAllianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(input.AgentID)
	if err != nil {
		return output, err
	}

	output.AgentAllianceList = agentAllianceList

	return output, nil
}

// アライアンス未申請がひとつでも存在することを確認する
type CheckAnyAllianceWithoutApplicationInput struct {
	MyAgentID        uint
	OtherAgentIDList []uint
}

type CheckAnyAllianceWithoutApplicationOutput struct {
	IDList []uint
}

func (i *AgentAllianceInteractorImpl) CheckAnyAllianceWithoutApplication(input CheckAnyAllianceWithoutApplicationInput) (CheckAnyAllianceWithoutApplicationOutput, error) {
	var (
		output CheckAnyAllianceWithoutApplicationOutput
	)

	// アライアンス申請の有無を確認
	agentAllianceList, err := i.agentAllianceRepository.GetByMyAgentIDAndOtherIDList(input.MyAgentID, input.OtherAgentIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 自分が申請済みのアライアンス申請リスト
	var appliedAllianceList []*entity.AgentAlliance

	// 自分が申請中かを確認
	// for _, alliance := range agentAllianceList {
	// 	if input.MyAgentID == alliance.Agent1ID && alliance.Agent1Request {
	// 		appliedAllianceList = append(appliedAllianceList, alliance)
	// 	} else if input.MyAgentID == alliance.Agent2ID && alliance.Agent2Request {
	// 		appliedAllianceList = append(appliedAllianceList, alliance)
	// 	}
	// }

	// アライアンス締結積みかを確認
	for _, alliance := range agentAllianceList {
		if alliance.Agent1Request && alliance.Agent2Request {
			appliedAllianceList = append(appliedAllianceList, alliance)
		}
	}

	// appliedAllianceListの値をマップに格納する
	existsOtherAgentID := make(map[uint]bool)
	for _, appliedAlliance := range appliedAllianceList {
		if input.MyAgentID == appliedAlliance.Agent1ID {
			existsOtherAgentID[appliedAlliance.Agent2ID] = true
		} else {
			existsOtherAgentID[appliedAlliance.Agent1ID] = true
		}
	}

	// input.OtherAgentIDListの各値がマップに存在しない場合、結果に追加する
	var unAppliedAgentIDList []uint
	for _, value := range input.OtherAgentIDList {
		if !existsOtherAgentID[value] {
			unAppliedAgentIDList = append(unAppliedAgentIDList, value)
		}
	}

	output.IDList = unAppliedAgentIDList

	return output, nil
}
