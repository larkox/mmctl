package commands

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mmctl/printer"

	"github.com/spf13/cobra"
)

func (s *MmctlUnitTestSuite) TestSearchUserCmd() {
	s.Run("Search for an existing user", func() {
		emailArg := "example@example.com"
		mockUser := model.User{Username: "ExampleUser", Email: emailArg}

		s.client.
			EXPECT().
			GetUserByEmail(emailArg, "").
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		err := searchUserCmdF(s.client, &cobra.Command{}, []string{emailArg})
		s.Require().Nil(err)
		s.Require().Equal(&mockUser, printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Search for a nonexistent user", func() {
		printer.Clean()
		arg := "example@example.com"

		s.client.
			EXPECT().
			GetUserByEmail(arg, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(arg, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetUser(arg, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		err := searchUserCmdF(s.client, &cobra.Command{}, []string{arg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Equal("Unable to find user 'example@example.com'", printer.GetErrorLines()[0])
	})
}

func (s *MmctlUnitTestSuite) TestUserInviteCmd() {
	s.Run("Invite user to an existing team by Id", func() {
		printer.Clean()
		argUser := "example@example.com"
		argTeam := "teamId"

		s.client.
			EXPECT().
			GetTeam(argTeam, "").
			Return(&model.Team{Id: argTeam}, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			InviteUsersToTeam(argTeam, []string{argUser}).
			Return(false, &model.Response{Error: nil}).
			Times(1)

		err := userInviteCmdF(s.client, &cobra.Command{}, []string{argUser, argTeam})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal("Invites may or may not have been sent.", printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Invite user to an existing team by name", func() {
		printer.Clean()
		argUser := "example@example.com"
		argTeam := "teamName"
		resultID := "teamId"

		s.client.
			EXPECT().
			GetTeam(argTeam, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeamByName(argTeam, "").
			Return(&model.Team{Id: resultID}, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			InviteUsersToTeam(resultID, []string{argUser}).
			Return(false, &model.Response{Error: nil}).
			Times(1)

		err := userInviteCmdF(s.client, &cobra.Command{}, []string{argUser, argTeam})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal("Invites may or may not have been sent.", printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Invite user to several existing teams by name and id", func() {
		printer.Clean()
		argUser := "example@example.com"
		argTeam := []string{"teamName1", "teamId2", "teamId3", "teamName4"}
		resultTeamModels := [4]*model.Team{
			&model.Team{Id: "teamId1"},
			&model.Team{Id: "teamId2"},
			&model.Team{Id: "teamId3"},
			&model.Team{Id: "teamId4"},
		}

		// Setup GetTeam
		s.client.
			EXPECT().
			GetTeam(argTeam[0], "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[1], "").
			Return(resultTeamModels[1], &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[2], "").
			Return(resultTeamModels[2], &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[3], "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		// Setup GetTeamByName
		s.client.
			EXPECT().
			GetTeamByName(argTeam[0], "").
			Return(resultTeamModels[0], &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeamByName(argTeam[3], "").
			Return(resultTeamModels[3], &model.Response{Error: nil}).
			Times(1)

		// Setup InvitUsersToTeam
		for _, resultTeamModel := range resultTeamModels {
			s.client.
				EXPECT().
				InviteUsersToTeam(resultTeamModel.Id, []string{argUser}).
				Return(false, &model.Response{Error: nil}).
				Times(1)
		}

		err := userInviteCmdF(s.client, &cobra.Command{}, append([]string{argUser}, argTeam...))
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), len(argTeam))
		for i := 0; i < len(argTeam); i++ {
			s.Require().Equal("Invites may or may not have been sent.", printer.GetLines()[i])
		}
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Invite user to an un-existing team", func() {
		printer.Clean()
		argUser := "example@example.com"
		argTeam := "unexistent"

		s.client.
			EXPECT().
			GetTeam(argTeam, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeamByName(argTeam, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		err := userInviteCmdF(s.client, &cobra.Command{}, []string{argUser, argTeam})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Len(printer.GetErrorLines(), 1)
		s.Require().Equal("Can't find team '"+argTeam+"'", printer.GetErrorLines()[0])
	})

	s.Run("Invite user to an existing team and fail invite", func() {
		printer.Clean()
		argUser := "example@example.com"
		argTeam := "teamId"
		resultName := "teamName"
		mockError := model.NewAppError("", "Mock Error", nil, "", 0)

		s.client.
			EXPECT().
			GetTeam(argTeam, "").
			Return(&model.Team{Id: argTeam, Name: resultName}, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			InviteUsersToTeam(argTeam, []string{argUser}).
			Return(false, &model.Response{Error: mockError}).
			Times(1)

		err := userInviteCmdF(s.client, &cobra.Command{}, []string{argUser, argTeam})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Len(printer.GetErrorLines(), 1)
		s.Require().Equal("Unable to invite user with email "+argUser+" to team "+resultName+". Error: "+mockError.Error(), printer.GetErrorLines()[0])
	})

	s.Run("Invite user to several existing and non-existing teams by name and id and reject one invite", func() {
		printer.Clean()
		argUser := "example@example.com"
		argTeam := []string{"teamName1", "unexistent", "teamId3", "teamName4", "reject", "teamId6"}
		resultTeamModels := [6]*model.Team{
			&model.Team{Id: "teamId1", Name: "teamName1"},
			nil,
			&model.Team{Id: "teamId3", Name: "teamName3"},
			&model.Team{Id: "teamId4", Name: "teamName4"},
			&model.Team{Id: "reject", Name: "rejectName"},
			&model.Team{Id: "teamId6", Name: "teamName6"},
		}
		mockError := model.NewAppError("", "Mock Error", nil, "", 0)

		// Setup GetTeam
		s.client.
			EXPECT().
			GetTeam(argTeam[0], "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[1], "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[2], "").
			Return(resultTeamModels[2], &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[3], "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[4], "").
			Return(resultTeamModels[4], &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeam(argTeam[5], "").
			Return(resultTeamModels[5], &model.Response{Error: nil}).
			Times(1)

		// Setup GetTeamByName
		s.client.
			EXPECT().
			GetTeamByName(argTeam[0], "").
			Return(resultTeamModels[0], &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeamByName(argTeam[1], "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetTeamByName(argTeam[3], "").
			Return(resultTeamModels[3], &model.Response{Error: nil}).
			Times(1)

		// Setup InvitUsersToTeam
		s.client.
			EXPECT().
			InviteUsersToTeam(resultTeamModels[0].Id, []string{argUser}).
			Return(false, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			InviteUsersToTeam(resultTeamModels[2].Id, []string{argUser}).
			Return(false, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			InviteUsersToTeam(resultTeamModels[3].Id, []string{argUser}).
			Return(false, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			InviteUsersToTeam(resultTeamModels[4].Id, []string{argUser}).
			Return(false, &model.Response{Error: mockError}).
			Times(1)

		s.client.
			EXPECT().
			InviteUsersToTeam(resultTeamModels[5].Id, []string{argUser}).
			Return(false, &model.Response{Error: nil}).
			Times(1)

		err := userInviteCmdF(s.client, &cobra.Command{}, append([]string{argUser}, argTeam...))
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 4)
		for i := 0; i < 4; i++ {
			s.Require().Equal("Invites may or may not have been sent.", printer.GetLines()[i])
		}
		s.Require().Len(printer.GetErrorLines(), 2)
		s.Require().Equal("Can't find team '"+argTeam[1]+"'", printer.GetErrorLines()[0])
		s.Require().Equal("Unable to invite user with email "+argUser+" to team "+resultTeamModels[4].Name+". Error: "+mockError.Error(), printer.GetErrorLines()[1])
	})
}

func (s *MmctlUnitTestSuite) TestUserCreateCmd() {
	mockUser := model.User{
		Username: "username",
		Password: "password",
		Email:    "email",
	}

	s.Run("Create user with email missing", func() {
		printer.Clean()

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Email is required: flag accessed but not defined: email", error.Error())
	})

	s.Run("Create user with username missing", func() {
		printer.Clean()

		command := cobra.Command{}
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Username is required: flag accessed but not defined: username", error.Error())
	})

	s.Run("Create user with password missing", func() {
		printer.Clean()

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Password is required: flag accessed but not defined: password", error.Error())
	})

	s.Run("Create a regular user", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Nil(error)
		s.Require().Equal(&mockUser, printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Create a regular user with client returning error", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: &model.AppError{Message: "Remote error"}}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Unable to create user. Error: : Remote error, ", error.Error())
	})

	s.Run("Create a sysAdmin user", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			UpdateUserRoles(mockUser.Id, "system_user system_admin").
			Return(true, &model.Response{Error: nil}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")
		command.Flags().Bool("system_admin", true, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Nil(error)
		s.Require().Equal(&mockUser, printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Create a sysAdmin user with client returning error", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			UpdateUserRoles(mockUser.Id, "system_user system_admin").
			Return(false, &model.Response{Error: &model.AppError{Message: "Remote error"}}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")
		command.Flags().Bool("system_admin", true, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Unable to update user roles. Error: : Remote error, ", error.Error())
	})
}
