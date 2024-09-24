//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/interfaces/handler"
	"github.com/spaceaiinc/autoscout-server/interfaces/repository"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

var wireSet = wire.NewSet(
	handler.WireSet,
	interactor.WireSet,
	repository.WireSet,
)

/**
	Handler
**/

// Seesion
func InitializeSessionHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SessionHandler) {
	wire.Build(wireSet)
	return
}

// Admin
func InitializeAdminHandler(db interfaces.SQLExecuter, appConfig config.App) (h handler.AdminHandler) {
	wire.Build(wireSet)
	return
}

// Agent
func InitializeAgentHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.AgentHandler) {
	wire.Build(wireSet)
	return
}

// AgentStaff
func InitializeAgentStaffHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.AgentStaffHandler) {
	wire.Build(wireSet)
	return
}

// AgentAlliance
func InitializeAgentAllianceHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.AgentAllianceHandler) {
	wire.Build(wireSet)
	return
}

// AgentRobot
func InitializeAgentRobotHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.AgentRobotHandler) {
	wire.Build(wireSet)
	return
}

// EnterpriseProfile
func InitializeEnterpriseProfileHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.EnterpriseProfileHandler) {
	wire.Build(wireSet)
	return
}

// BillingAddress
func InitializeBillingAddressHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.BillingAddressHandler) {
	wire.Build(wireSet)
	return
}

// JobInformation
func InitializeJobInformationHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.JobInformationHandler) {
	wire.Build(wireSet)
	return
}

// JobSeeker
func InitializeJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) (h handler.JobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// Task
func InitializeTaskHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.TaskHandler) {
	wire.Build(wireSet)
	return
}

// MessageTemplate
func InitializeMessageTemplateHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.MessageTemplateHandler) {
	wire.Build(wireSet)
	return
}

// SelectionQuestionnaire
func InitializeSelectionQuestionnaireHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SelectionQuestionnaireHandler) {
	wire.Build(wireSet)
	return
}

// ChatGroupWithJobSeeker
func InitializeChatGroupWithJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.ChatGroupWithJobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// ChatMessageWithJobSeeker
func InitializeChatMessageWithJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.ChatMessageWithJobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// ChatGroupWithAgent
func InitializeChatGroupWithAgentHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.ChatGroupWithAgentHandler) {
	wire.Build(wireSet)
	return
}

// ChatThreadWithAgent
func InitializeChatThreadWithAgentHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.ChatThreadWithAgentHandler) {
	wire.Build(wireSet)
	return
}

// ChatMessageWithAgent
func InitializeChatMessageWithAgentHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.ChatMessageWithAgentHandler) {
	wire.Build(wireSet)
	return
}

// InterviewTask
func InitializeInterviewTaskHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.InterviewTaskHandler) {
	wire.Build(wireSet)
	return
}

// InterviewAdjustmentTemplate
func InitializeInterviewAdjustmentTemplateHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.InterviewAdjustmentTemplateHandler) {
	wire.Build(wireSet)
	return
}

// ScoutService
func InitializeScoutServiceHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slackAPI config.Slack) (h handler.ScoutServiceHandler) {
	wire.Build(wireSet)
	return
}

// Sale
func InitializeSaleHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.SaleHandler) {
	wire.Build(wireSet)
	return
}

// AgentMonthlySale
func InitializeAgentMonthlySaleHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.AgentMonthlySaleHandler) {
	wire.Build(wireSet)
	return
}

// AgentStaffMonthlySale
func InitializeAgentStaffMonthlySaleHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.AgentStaffMonthlySaleHandler) {
	wire.Build(wireSet)
	return
}

// Schedule
func InitializeScheduleHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.ScheduleHandler) {
	wire.Build(wireSet)
	return
}

// Dashboard
func InitializeDashboardHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.DashboardHandler) {
	wire.Build(wireSet)
	return
}

// JobSeekerSchedule
func InitializeJobSeekerScheduleHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.JobSeekerScheduleHandler) {
	wire.Build(wireSet)
	return
}

// InitialEnterpriseImporter
func InitializeInitialEnterpriseImporterHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.InitialEnterpriseImporterHandler) {
	wire.Build(wireSet)
	return
}

// NotificationForUser
func InitializeNotificationForUserHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.NotificationForUserHandler) {
	wire.Build(wireSet)
	return
}

// EmailWithJobSeeker
func InitializeEmailWithJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.EmailWithJobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// DeploymentReflection
func InitializeDeploymentReflectionHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.DeploymentReflectionHandler) {
	wire.Build(wireSet)
	return
}

// AgentInflowChannelOption
func InitializeAgentInflowChannelOptionHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.AgentInflowChannelOptionHandler) {
	wire.Build(wireSet)
	return
}

// SendingCustomer
func InitializeSendingCustomerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.SendingCustomerHandler) {
	wire.Build(wireSet)
	return
}

// SendingJobSeeker
func InitializeSendingJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.SendingJobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// ChatGroupWithSendingJobSeeker
func InitializeChatGroupWithSendingJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.ChatGroupWithSendingJobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// ChatMessageWithSendingJobSeeker
func InitializeChatMessageWithSendingJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.ChatMessageWithSendingJobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// EmailWithSendingJobSeeker
func InitializeEmailWithSendingJobSeekerHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (h handler.EmailWithSendingJobSeekerHandler) {
	wire.Build(wireSet)
	return
}

// SendingEnterprise
func InitializeSendingEnterpriseHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SendingEnterpriseHandler) {
	wire.Build(wireSet)
	return
}

// SendingBillingAddress
func InitializeSendingBillingAddressHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SendingBillingAddressHandler) {
	wire.Build(wireSet)
	return
}

// SendingJobInformation
func InitializeSendingJobInformationHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SendingJobInformationHandler) {
	wire.Build(wireSet)
	return
}

// SendingPhase
func InitializeSendingPhaseHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SendingPhaseHandler) {
	wire.Build(wireSet)
	return
}

// SendingJobSeekerDesiredJobInformation
func InitializeSendingJobSeekerDesiredJobInformationHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SendingJobSeekerDesiredJobInformationHandler) {
	wire.Build(wireSet)
	return
}

// SendingSale
func InitializeSendingSaleHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (h handler.SendingSaleHandler) {
	wire.Build(wireSet)
	return
}

// GoogleAuthentication
func InitializeGoogleAuthenticationHandler(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, googleAPI config.GoogleAPI) (h handler.GoogleAuthenticationHandler) {
	wire.Build(wireSet)
	return
}

/**
	Interactor
**/

// Session
func InitializeSessionInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SessionInteractor) {
	wire.Build(wireSet)
	return
}

// Admin
func InitializeAdminInteractor(db interfaces.SQLExecuter, appConfig config.App) (i interactor.AdminInteractor) {
	wire.Build(wireSet)
	return
}

// Agent
func InitializeAgentInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.AgentInteractor) {
	wire.Build(wireSet)
	return
}

// AgentStaff
func InitializeAgentStaffInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.AgentStaffInteractor) {
	wire.Build(wireSet)
	return
}

// AgentAlliance
func InitializeAgentAllianceInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.AgentAllianceInteractor) {
	wire.Build(wireSet)
	return
}

// AgentRobot
func InitializeAgentRobotInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.AgentRobotInteractor) {
	wire.Build(wireSet)
	return
}

// EnterpriseProfile
func InitializeEnterpriseProfileInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.EnterpriseProfileInteractor) {
	wire.Build(wireSet)
	return
}

// BillingAddress
func InitializeBillingAddressInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.BillingAddressInteractor) {
	wire.Build(wireSet)
	return
}

// JobInformation
func InitializeJobInformationInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.JobInformationInteractor) {
	wire.Build(wireSet)
	return
}

// JobSeeker
func InitializeJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) (i interactor.JobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// Task
func InitializeTaskInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.TaskInteractor) {
	wire.Build(wireSet)
	return
}

// SelectionQuestionnaire
func InitializeSelectionQuestionnaireInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SelectionQuestionnaireInteractor) {
	wire.Build(wireSet)
	return
}

// ChatGroupWithJobSeeker
func InitializeChatGroupWithJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.ChatGroupWithJobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// ChatMessageWithJobSeeker
func InitializeChatMessageWithJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.ChatMessageWithJobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// ChatGroupWithAgent
func InitializeChatGroupWithAgentInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.ChatGroupWithAgentInteractor) {
	wire.Build(wireSet)
	return
}

// ChatThreadWithAgent
func InitializeChatThreadWithAgentInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.ChatThreadWithAgentInteractor) {
	wire.Build(wireSet)
	return
}

// ChatMessageWithAgent
func InitializeChatMessageWithAgentInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.ChatMessageWithAgentInteractor) {
	wire.Build(wireSet)
	return
}

// InterviewTask
func InitializeInterviewTaskInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.InterviewTaskInteractor) {
	wire.Build(wireSet)
	return
}

// InterviewAdjustmentTemplate
func InitializeInterviewAdjustmentTemplateInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.InterviewAdjustmentTemplateInteractor) {
	wire.Build(wireSet)
	return
}

// ScoutService
func InitializeScoutServiceInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slackAPI config.Slack) (i interactor.ScoutServiceInteractor) {
	wire.Build(wireSet)
	return
}

// Sale
func InitializeSaleInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.SaleInteractor) {
	wire.Build(wireSet)
	return
}

// AgentMonthlySale
func InitializeAgentMonthlySaleInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.AgentMonthlySaleInteractor) {
	wire.Build(wireSet)
	return
}

// AgentStaffMonthlySale
func InitializeAgentStaffMonthlySaleInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.AgentStaffMonthlySaleInteractor) {
	wire.Build(wireSet)
	return
}

// Schedule
func InitializeScheduleInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.ScheduleInteractor) {
	wire.Build(wireSet)
	return
}

// Dashboard
func InitializeDashboardInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.DashboardInteractor) {
	wire.Build(wireSet)
	return
}

// JobSeekerSchedule
func InitializeJobSeekerScheduleInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.JobSeekerScheduleInteractor) {
	wire.Build(wireSet)
	return
}

// InitialEnterpriseImporter
func InitializeInitialEnterpriseImporterInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.InitialEnterpriseImporterInteractor) {
	wire.Build(wireSet)
	return
}

// NotificationForUser
func InitializeNotificationForUserInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.NotificationForUserInteractor) {
	wire.Build(wireSet)
	return
}

// EmailWithJobSeeker
func InitializeEmailWithJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.EmailWithJobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// DeploymentReflection
func InitializeDeploymentReflectionInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.DeploymentReflectionInteractor) {
	wire.Build(wireSet)
	return
}

// AgentInflowChannelOption
func InitializeAgentInflowChannelOptionInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.AgentInflowChannelOptionInteractor) {
	wire.Build(wireSet)
	return
}

// SendingCustomer
func InitializeSendingCustomerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.SendingCustomerInteractor) {
	wire.Build(wireSet)
	return
}

// SendingJobSeeker
func InitializeSendingJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.SendingJobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// ChatGroupWithSendingJobSeeker
func InitializeChatGroupWithSendingJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.ChatGroupWithSendingJobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// ChatMessageWithSendingJobSeeker
func InitializeChatMessageWithSendingJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.ChatMessageWithSendingJobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// EmailWithSendingJobSeeker
func InitializeEmailWithSendingJobSeekerInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, oneSignal config.OneSignal) (i interactor.EmailWithSendingJobSeekerInteractor) {
	wire.Build(wireSet)
	return
}

// SendingEnterprise
func InitializeSendingEnterpriseInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SendingEnterpriseInteractor) {
	wire.Build(wireSet)
	return
}

// SendingBillingAddress
func InitializeSendingBillingAddressInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SendingBillingAddressInteractor) {
	wire.Build(wireSet)
	return
}

// SendingJobInformation
func InitializeSendingJobInformationInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SendingJobInformationInteractor) {
	wire.Build(wireSet)
	return
}

// SendingPhase
func InitializeSendingPhaseInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SendingPhaseInteractor) {
	wire.Build(wireSet)
	return
}

// SendingJobSeekerDesiredJobInformation
func InitializeSendingJobSeekerDesiredJobInformationInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SendingJobSeekerDesiredJobInformationInteractor) {
	wire.Build(wireSet)
	return
}

// SendingSale
func InitializeSendingSaleInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid) (i interactor.SendingSaleInteractor) {
	wire.Build(wireSet)
	return
}

// GoogleAuthentication
func InitializeGoogleAuthenticationInteractor(fb usecase.Firebase, db interfaces.SQLExecuter, sendgrid config.Sendgrid, googleAPI config.GoogleAPI) (i interactor.GoogleAuthenticationInteractor) {
	wire.Build(wireSet)
	return
}
