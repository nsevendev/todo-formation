package task

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"todof/internal/auth"
	initializer "todof/internal/init"
	"todof/testsetup"
)

var (
	ctx         = context.Background()
	serviceTask TaskServiceInterface
	serviceUser auth.UserServiceInterface
	ctxCanceled context.Context
	cancelFunc  context.CancelFunc
	userDto     auth.UserCreateDto
	taskDto     TaskCreateDto
)

func TestMain(m *testing.M) {
	userDto = auth.CreateDtoFaker()
	taskDto = CreateDtoFaker()
	userRepo := auth.NewUserRepo(initializer.Db)
	serviceUser = auth.NewUserService(userRepo, "jwttest")
	serviceTask = NewTaskService(NewTaskRepo(initializer.Db), userRepo)
	ctxCanceled, cancelFunc = context.WithCancel(ctx)
	cancelFunc()

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
	code := m.Run()
	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")

	os.Exit(code)
}

func TestCreate_OneTaskCreated(t *testing.T) {
	var (
		err         error
		userCreated *auth.User
		taskCreated *Task
	)
	testsetup.LogNameTest(t, "one Tache created for one user")

	userCreated, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	taskCreated, err = serviceTask.Create(ctx, taskDto, userCreated.ID)
	testsetup.IsNull(t, err)
	testsetup.Equal(t, taskCreated.Label, taskDto.Label)
	testsetup.Equal(t, taskCreated.IdUser, userCreated.ID)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestCreate_ErrorDatabase(t *testing.T) {
	var (
		err  error
		user *auth.User
	)

	testsetup.LogNameTest(t, "create task with context canceled")

	user, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	_, err = serviceTask.Create(ctxCanceled, taskDto, user.ID)
	testsetup.IsNotNull(t, err)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestGetAllByUser_GetAllTaskByUser(t *testing.T) {
	var (
		err   error
		user  *auth.User
		task1 *Task
		task2 *Task
		tasks []Task
	)

	testsetup.LogNameTest(t, "get all tasks by user")

	user, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	task1, err = serviceTask.Create(ctx, taskDto, user.ID)
	testsetup.IsNull(t, err)
	task2, err = serviceTask.Create(ctx, taskDto, user.ID)
	testsetup.IsNull(t, err)
	tasks, err = serviceTask.GetAllByUser(ctx, user.ID)
	testsetup.IsNull(t, err)

	testsetup.Equal(t, len(tasks), 2)
	testsetup.Equal(t, tasks[0].ID, task1.ID)
	testsetup.Equal(t, tasks[1].ID, task2.ID)
	testsetup.Equal(t, tasks[0].Label, task1.Label)
	testsetup.Equal(t, tasks[1].Label, task2.Label)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestGetAllByUser_GetAllTaskByUserWithUserNotExist(t *testing.T) {
	var (
		err   error
		user  *auth.User
		tasks []Task
	)

	testsetup.LogNameTest(t, "get all tasks by user with user not exist")

	user, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	_, err = serviceTask.Create(ctx, taskDto, user.ID)
	testsetup.IsNull(t, err)
	_, err = serviceTask.Create(ctx, taskDto, user.ID)
	testsetup.IsNull(t, err)
	tasks, err = serviceTask.GetAllByUser(ctx, primitive.NewObjectID())
	testsetup.IsNull(t, err)
	testsetup.IsEmptySlice(t, tasks)
	testsetup.Equal(t, len(tasks), 0)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestGetAllByUser_ErrorDatabase(t *testing.T) {
	var (
		err  error
		user *auth.User
	)

	testsetup.LogNameTest(t, "get all tasks by user with context canceled")

	user, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	_, err = serviceTask.Create(ctx, taskDto, user.ID)
	testsetup.IsNull(t, err)
	_, err = serviceTask.Create(ctx, taskDto, user.ID)
	testsetup.IsNull(t, err)
	_, err = serviceTask.GetAllByUser(ctxCanceled, user.ID)
	testsetup.IsNotNull(t, err)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestGetAllByUser_TaskIsNotValidate(t *testing.T) {
	var (
		err   error
		user  *auth.User
		tasks []Task
	)

	testsetup.LogNameTest(t, "get all tasks by user with task not validate")

	user, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)

	// on insère un document mal formé
	_, err = initializer.Db.Collection("tasks").InsertOne(ctx, bson.M{
		"id_user": user.ID,
		"label":   1234,
		"done":    false,
	}, options.InsertOne().SetBypassDocumentValidation(true))
	testsetup.IsNull(t, err)

	tasks, err = serviceTask.GetAllByUser(ctx, user.ID)
	testsetup.IsEmptySlice(t, tasks)
	testsetup.IsNotNull(t, err)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestUpdateOneDonePropertyByUser_Success(t *testing.T) {
	var (
		err         error
		userCreated *auth.User
		taskCreated *Task
		tasks       []Task
	)

	testsetup.LogNameTest(t, "update one done property by user success")

	userCreated, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	taskCreated, err = serviceTask.Create(ctx, taskDto, userCreated.ID)
	testsetup.IsNull(t, err)
	testsetup.Equal(t, taskCreated.Done, false)

	err = serviceTask.UpdateOneDonePropertyByUser(ctx, userCreated.ID, taskCreated.ID)
	testsetup.IsNull(t, err)
	tasks, err = serviceTask.GetAllByUser(ctx, userCreated.ID)
	testsetup.IsNull(t, err)
	testsetup.Equal(t, len(tasks), 1)
	testsetup.Equal(t, tasks[0].ID, taskCreated.ID)
	testsetup.Equal(t, tasks[0].Done, true)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestUpdateOneDonePropertyByUser_WithTaskNotExist(t *testing.T) {
	var (
		err         error
		userCreated *auth.User
	)

	testsetup.LogNameTest(t, "update one done property by user with task not exist")

	userCreated, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	err = serviceTask.UpdateOneDonePropertyByUser(ctx, userCreated.ID, primitive.NewObjectID())
	testsetup.IsNotNull(t, err)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestUpdateOneDonePropertyByUser_WithUserNotExist(t *testing.T) {
	var (
		err         error
		userCreated *auth.User
		taskCreated *Task
	)

	testsetup.LogNameTest(t, "update one done property by user with user not exist")

	userCreated, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	taskCreated, err = serviceTask.Create(ctx, taskDto, userCreated.ID)
	testsetup.IsNull(t, err)

	err = serviceTask.UpdateOneDonePropertyByUser(ctx, primitive.NewObjectID(), taskCreated.ID)
	testsetup.IsNotNull(t, err)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestUpdateOneDonePropertyByUser_WithUserExistButHasNotTask(t *testing.T) {
	var (
		err          error
		userCreated  *auth.User
		userCreated2 *auth.User
		taskCreated  *Task
	)

	testsetup.LogNameTest(t, "update one done property by user with user exist but has not task")

	userCreated, err = serviceUser.Register(ctx, userDto)
	testsetup.IsNull(t, err)
	userDto2 := auth.CreateDtoFaker()
	userCreated2, err = serviceUser.Register(ctx, userDto2)
	testsetup.IsNull(t, err)
	taskCreated, err = serviceTask.Create(ctx, taskDto, userCreated.ID)
	testsetup.IsNull(t, err)

	err = serviceTask.UpdateOneDonePropertyByUser(ctx, userCreated2.ID, taskCreated.ID)
	testsetup.IsNotNull(t, err)

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

/*
func TestUpdateOneDonePropertyByUser(t *testing.T) {
	tests := []struct {
		name   string
		idUser primitive.ObjectID
		idTask primitive.ObjectID
		isErr  bool
	}{
		{"test success", usersIds[0], tasksIds[0], false},
		{"test echec avec task introuvale", usersIds[0], primitive.NewObjectID(), true},
		{"test echec mongo", usersIds[0], tasksIds[0], true},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			err := s.UpdateOneDonePropertyByUser(cancelCtx, tt.idUser, tt.idTask)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			err := s.UpdateOneDonePropertyByUser(ctx, tt.idUser, tt.idTask)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}
	}
}

func TestUpdateOneLabelPropertyByUser(t *testing.T) {
	tests := []struct {
		name   string
		idUser primitive.ObjectID
		idTask primitive.ObjectID
		label  string
		isErr  bool
	}{
		{"test success", usersIds[0], tasksIds[0], "label updated", false},
		{"test echec mongo", usersIds[0], tasksIds[0], "label updated", true},
	}

	for _, tt := range tests {
		updateDto := TaskUpdateLabelDto{
			Label: tt.label,
		}

		switch tt.name {
		case "test echec mongo":
			err := s.UpdateOneLabelPropertyByUser(cancelCtx, tt.idUser, tt.idTask, updateDto)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			err := s.UpdateOneLabelPropertyByUser(ctx, tt.idUser, tt.idTask, updateDto)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}
	}
}

func TestDeleteOneByUser(t *testing.T) {
	tests := []struct {
		name   string
		idUser primitive.ObjectID
		idTask primitive.ObjectID
		isErr  bool
	}{
		{"test success", usersIds[0], tasksIds[0], false},
		{"test echec mongodb", usersIds[0], tasksIds[0], true},
		{"test aucune tâche supprimée", usersIds[0], primitive.NewObjectID(), true},
	}

	for _, tt := range tests {
		if tt.name == "test echec mongodb" {
			err := s.DeleteOneByUser(cancelCtx, tt.idUser, tt.idTask)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		} else {
			err := s.DeleteOneByUser(ctx, tt.idUser, tt.idTask)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}
	}
}

func TestDeleteManyByUser(t *testing.T) {

	tests := []struct {
		name  string
		setup func() (userID primitive.ObjectID, taskID primitive.ObjectID)
		isErr bool
	}{
		{
			name: "test success",
			setup: func() (primitive.ObjectID, primitive.ObjectID) {
				userID := setupCreateUser(t, "taskTest3@gmail.com")
				taskID := setupCreateTask(t, userID, "test label")
				return userID, taskID
			},
			isErr: false,
		},
		{
			name: "test echec mongo",
			setup: func() (primitive.ObjectID, primitive.ObjectID) {
				userID := setupCreateUser(t, "taskTest4@gmail.com")
				taskID := setupCreateTask(t, userID, "test label")
				return userID, taskID
			},
			isErr: true,
		},
		{
			name: "test aucune tâche supprimée",
			setup: func() (primitive.ObjectID, primitive.ObjectID) {
				userID := setupCreateUser(t, "taskTest5@gmail.com")
				taskID := setupCreateTask(t, userID, "test label")
				return userID, taskID
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		userID, taskID := tt.setup()

		switch tt.name {
		case "test echec mongo":
			err := s.DeleteManyByUser(cancelCtx, userID, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		case "test aucune tâche supprimée":
			err := s.DeleteManyByUser(ctx, userID, []primitive.ObjectID{primitive.NewObjectID()})

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			err := s.DeleteManyByUser(ctx, userID, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("got error %v, expected error: %v", err, tt.isErr)
			}
		}
	}
}

func TestDeleteById(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (userID primitive.ObjectID, taskID primitive.ObjectID)
		isErr bool
	}{
		{
			name: "test success",
			setup: func() (primitive.ObjectID, primitive.ObjectID) {
				userID := setupCreateUser(t, "taskTest6@gmail.com")
				taskID := setupCreateTask(t, userID, "test label")
				return userID, taskID
			},
			isErr: false,
		},
		{
			name: "test echec mongo",
			setup: func() (primitive.ObjectID, primitive.ObjectID) {
				userID := setupCreateUser(t, "taskTest7@gmail.com")
				taskID := setupCreateTask(t, userID, "test label")
				return userID, taskID
			},
			isErr: true,
		},
		{
			name: "test aucune tâche supprimée",
			setup: func() (primitive.ObjectID, primitive.ObjectID) {
				userID := setupCreateUser(t, "taskTest8@gmail.com")
				taskID := setupCreateTask(t, userID, "test label")
				return userID, taskID
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		_, taskID := tt.setup()

		switch tt.name {
		case "test echec mongo":
			err := s.DeleteById(cancelCtx, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		case "test aucune tâche supprimée":
			err := s.DeleteById(ctx, []primitive.ObjectID{primitive.NewObjectID()})

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			err := s.DeleteById(ctx, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("got error %v, expected error: %v", err, tt.isErr)
			}
		}
	}
}

func TestDeleteAllTasks(t *testing.T) {
	tests := []struct {
		name  string
		isErr bool
	}{
		{"test success", false},
		{"test echec mongodb", true},
	}

	for _, tt := range tests {

		if tt.name == "test echec mongodb" {
			err := s.DeleteAllTasks(cancelCtx)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		} else {
			err := s.DeleteAllTasks(ctx)

			if (err != nil) != tt.isErr {
				t.Errorf("got error %v, expected error: %v", err, tt.isErr)
			}
		}
	}
} */
