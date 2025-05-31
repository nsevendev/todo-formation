package task

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
	"testing"
	"todof/internal/auth"
	initializer "todof/internal/init"
	"todof/testsetup"
)

var (
	service     TaskServiceInterface
	userService auth.UserServiceInterface
	ctx         = context.Background()
	ctxCanceled context.Context
	cancelFunc  context.CancelFunc
	userDto     auth.UserCreateDto
	taskDto     TaskCreateDto
)

func TestMain(m *testing.M) {
	userRepo := auth.NewUserRepo(initializer.Db)
	userService = auth.NewUserService(userRepo, "jwttest")
	service = NewTaskService(NewTaskRepo(initializer.Db), userRepo)
	ctxCanceled, cancelFunc = context.WithCancel(ctx)
	cancelFunc()

	userDto = auth.UserCreateDto{
		Email:    "taskTest@gmail.com",
		Password: "123456",
		Username: "taskTest",
	}

	taskDto = TaskCreateDto{
		Label: "task test",
	}

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
	code := m.Run()
	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")

	os.Exit(code)
}

func TestCreate_Success(t *testing.T) {
	var (
		err  error
		user *auth.User
		task *Task
	)

	user, err = userService.Register(ctx, userDto)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	task, err = service.Create(ctx, taskDto, user.ID)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	if task.Label != taskDto.Label {
		testsetup.Except(t, task.Label, taskDto.Label)
	}

	testsetup.Success(t, "label task => "+task.Label, "label taskdto => "+taskDto.Label)
	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestCreate_MongoCtxCanceled(t *testing.T) {
	var (
		err  error
		user *auth.User
	)

	user, err = userService.Register(ctx, userDto)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	_, err = service.Create(ctxCanceled, taskDto, user.ID)
	// on veut une erreur donc on teste que err n'est pas nil
	if err == nil {
		testsetup.ErrorSuccess(t)
	} else {
		testsetup.Success(t, "context canceled", err.Error())
	}

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestGetAllByUser_Success(t *testing.T) {
	var (
		err   error
		user  *auth.User
		task1 *Task
		task2 *Task
		tasks []Task
	)
	excepted := 2

	user, err = userService.Register(ctx, userDto)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	task1, err = service.Create(ctx, taskDto, user.ID)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	task2, err = service.Create(ctx, taskDto, user.ID)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	tasks, err = service.GetAllByUser(ctx, user.ID)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	if len(tasks) != excepted {
		testsetup.Except(t, strconv.Itoa(len(tasks)), strconv.Itoa(excepted))
	}

	if tasks[0].ID != task1.ID {
		testsetup.Except(t, tasks[0].ID.Hex(), task1.ID.Hex())
	}

	if tasks[1].ID != task2.ID {
		testsetup.Except(t, tasks[1].ID.Hex(), task2.ID.Hex())
	}

	if tasks[0].Label != task1.Label {
		testsetup.Except(t, tasks[0].Label, task1.Label)
	}

	if tasks[1].Label != task2.Label {
		testsetup.Except(t, tasks[1].Label, task2.Label)
	}

	testsetup.Success(t, "len(tasks) => "+strconv.Itoa(len(tasks)), "len(tasks) => "+strconv.Itoa(excepted))
	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestGetAllByUser_MongoCtxCanceld(t *testing.T) {
	var (
		err  error
		user *auth.User
	)

	user, err = userService.Register(ctx, userDto)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	_, err = service.Create(ctx, taskDto, user.ID)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	_, err = service.Create(ctx, taskDto, user.ID)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	_, err = service.GetAllByUser(ctxCanceled, user.ID)
	// on veut une erreur donc on teste que err n'est pas nil
	if err == nil {
		testsetup.ErrorSuccess(t)
	} else {
		testsetup.Success(t, "context canceled", err.Error())
	}

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

func TestGetAllByUser_TaskIsNotValidate(t *testing.T) {
	var (
		err   error
		user  *auth.User
		tasks []Task
	)

	user, err = userService.Register(ctx, userDto)
	if err != nil {
		testsetup.Error(t, err.Error())
	}

	// on insère un document mal formé
	_, err = initializer.Db.Collection("tasks").InsertOne(ctx, bson.M{
		"id_user": user.ID,
		"label":   1234,
		"done":    false,
	}, options.InsertOne().SetBypassDocumentValidation(true))
	if err != nil {
		testsetup.Error(t, "Erreur lors de l'insertion du document mal formé : "+err.Error())
	}

	tasks, err = service.GetAllByUser(ctx, user.ID)
	if tasks != nil {
		testsetup.Except(t, "nil", "not nil")
	}

	// on veut une erreur donc on teste que err n'est pas nil
	if err == nil {
		testsetup.ErrorSuccess(t)
	} else {
		testsetup.Success(t, "error decoding key label: cannot decode 32-bit integer into a string type", err.Error())
	}

	testsetup.CleanCollections(ctx, initializer.Db, "tasks", "users")
}

/*
func TestGetAllByUser(t *testing.T) {

	tests := []struct {
		name   string
		userId primitive.ObjectID
		isTask bool
		isErr  bool
	}{
		{"test success", usersIds[0], true, false},
		{"test avec user sans task", setupCreateUser(t, "taskTest2@gmail.com"), true, false},
		{"test echec mongo", usersIds[0], false, true},
		{"test document mal formé", usersIds[0], false, true},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			_, err := s.GetAllByUser(cancelCtx, tt.userId)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		case "test document mal formé":
			_, err := taskCollection.InsertOne(ctx, bson.M{
				"id_user": tt.userId,
				"label":   1234,
				"done":    false,
			}, options.InsertOne().SetBypassDocumentValidation(true))
			if err != nil {
				t.Fatalf("Erreur lors de l'insertion du document mal formé: %v", err)
			}

			_, err = s.GetAllByUser(ctx, tt.userId)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			task, err := s.GetAllByUser(ctx, tt.userId)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
			if (task == nil) == tt.isTask {
				t.Errorf("%s: got task %v, expect task %v", tt.name, task, tt.isTask)
			}
		}
	}
}

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
