package service_test

import (
	"context"
	"errors"
	"optionhub-service/internal/config"
	"optionhub-service/internal/model"
	"optionhub-service/internal/service"
	"testing"

	"github.com/golang/mock/gomock"
	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_AddOS(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		osName := "ubuntu"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddOS(gomock.Any(), osName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddOs(ctx, &optionhubproto.AddIn{Value: osName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: osName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		osName := "macOS"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := service.NewMockDBRepo(ctrl)

		s := service.NewService(mockRepo)

		_, err := s.AddOs(ctx, &optionhubproto.AddIn{Value: osName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "cannot find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		osName := "windows"

		var expectedID int64

		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddOS(gomock.Any(), osName, uuid).Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddOs(ctx, &optionhubproto.AddIn{Value: osName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetOsById(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		expectedOsName := "ubuntu"

		var id int64 = 3

		mockRepo.EXPECT().GetOsByID(gomock.Any(), id).Return(expectedOsName, nil)

		s := service.NewService(mockRepo)
		osName, err := s.GetOsByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, osName, &optionhubproto.GetByIdOut{Id: id, Value: expectedOsName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetOsByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetOsByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_GetOsBySearchName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		expectedNames := []model.CategoryItem{
			{ID: 1, Name: "ubuntu"},
			{ID: 2, Name: "ubuntuu"},
			{ID: 5, Name: "ubububu"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "ubuntu"},
				{Id: 2, Name: "ubuntuu"},
				{Id: 5, Name: "ubububu"},
			},
		}
		search := "ub"

		mockRepo.EXPECT().GetOsBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		osNames, err := s.GetOsBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, osNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		search := "w"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Name: "windows"},
			{ID: 2, Name: "wsl"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "windows"},
				{Id: 2, Name: "wsl"},
			},
		}

		mockRepo.EXPECT().GetOsPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		osNames, err := s.GetOsBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, osNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		search := "wi"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetOsBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetOsBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetAllOs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_all_os_ok", func(t *testing.T) {
		expectedNames := []model.CategoryItem{
			{ID: 1, Name: "ubuntu"},
			{ID: 2, Name: "Mac OS"},
			{ID: 5, Name: "Windows"},
		}
		expectedRes := &optionhubproto.GetAllOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "ubuntu"},
				{Id: 2, Name: "Mac OS"},
				{Id: 5, Name: "Windows"},
			},
		}

		mockRepo.EXPECT().GetAllOs().Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		osNames, err := s.GetAllOs(ctx, &optionhubproto.EmptyOptionhub{})
		assert.NoError(t, err)
		assert.Equal(t, osNames, expectedRes)
	})

	t.Run("get_all_os_err", func(t *testing.T) {
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetAllOs().Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetAllOs(ctx, &optionhubproto.EmptyOptionhub{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetWorkPlaceBySearchName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		expectedNames := []model.CategoryItem{
			{ID: 1, Name: "avito"},
			{ID: 2, Name: "avitoo"},
			{ID: 5, Name: "avivito"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "avito"},
				{Id: 2, Name: "avitoo"},
				{Id: 5, Name: "avivito"},
			},
		}
		search := "av"

		mockRepo.EXPECT().GetWorkPlaceBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		workPlaceNames, err := s.GetWorkPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, workPlaceNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		search := "w"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Name: "wildberries"},
			{ID: 2, Name: "ozon"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "wildberries"},
				{Id: 2, Name: "ozon"},
			},
		}

		mockRepo.EXPECT().GetWorkPlacePreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		workPlaceNames, err := s.GetWorkPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, workPlaceNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		search := "wi"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetWorkPlaceBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetWorkPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetWorkPlaceById(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		expectedWorkPlaceName := "avito"

		var id int64 = 3

		mockRepo.EXPECT().GetWorkPlaceByID(gomock.Any(), id).Return(expectedWorkPlaceName, nil)

		s := service.NewService(mockRepo)
		workPlaceName, err := s.GetWorkPlaceById(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, workPlaceName, &optionhubproto.GetByIdOut{Id: id, Value: expectedWorkPlaceName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetWorkPlaceByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetWorkPlaceById(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddWorkPlace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		workPlaceName := "avito"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddWorkPlace(gomock.Any(), workPlaceName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddWorkPlace(ctx, &optionhubproto.AddIn{Value: workPlaceName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: workPlaceName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		workPlaceName := "wildberries"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := service.NewMockDBRepo(ctrl)

		s := service.NewService(mockRepo)

		_, err := s.AddWorkPlace(ctx, &optionhubproto.AddIn{Value: workPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		workPlaceName := "ozon"

		var expectedID int64

		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddWorkPlace(gomock.Any(), workPlaceName, uuid).Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddWorkPlace(ctx, &optionhubproto.AddIn{Value: workPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetStudyPlaceBySearchName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		expectedNames := []model.CategoryItem{
			{ID: 1, Name: "FU"},
			{ID: 2, Name: "HSE"},
			{ID: 5, Name: "MGIMO"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "FU"},
				{Id: 2, Name: "HSE"},
				{Id: 5, Name: "MGIMO"},
			},
		}
		search := "HS"

		mockRepo.EXPECT().GetStudyPlaceBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		studyPlaceNames, err := s.GetStudyPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, studyPlaceNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		search := "m"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Name: "mirea"},
			{ID: 2, Name: "mgimo"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "mirea"},
				{Id: 2, Name: "mgimo"},
			},
		}

		mockRepo.EXPECT().GetStudyPlacePreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		studyPlaceNames, err := s.GetStudyPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, studyPlaceNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		search := "hs"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetStudyPlaceBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetStudyPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetStudyPlaceById(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		expectedStudyPlaceName := "MAI"

		var id int64 = 3

		mockRepo.EXPECT().GetStudyPlaceByID(gomock.Any(), id).Return(expectedStudyPlaceName, nil)

		s := service.NewService(mockRepo)
		studyPlaceName, err := s.GetStudyPlaceById(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, studyPlaceName, &optionhubproto.GetByIdOut{Id: id, Value: expectedStudyPlaceName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetStudyPlaceByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetStudyPlaceById(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddStudyPlace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		studyPlaceName := "rudn"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddStudyPlace(gomock.Any(), studyPlaceName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddStudyPlace(ctx, &optionhubproto.AddIn{Value: studyPlaceName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: studyPlaceName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		studyPlaceName := "rudn"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := service.NewMockDBRepo(ctrl)

		s := service.NewService(mockRepo)

		_, err := s.AddStudyPlace(ctx, &optionhubproto.AddIn{Value: studyPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		studyPlaceName := "rudn"

		var expectedID int64

		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddStudyPlace(gomock.Any(), studyPlaceName, uuid).Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddStudyPlace(ctx, &optionhubproto.AddIn{Value: studyPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetHobbyBySearchName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		expectedNames := []model.CategoryItem{
			{ID: 1, Name: "driving"},
			{ID: 2, Name: "painting"},
			{ID: 5, Name: "swimming"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "driving"},
				{Id: 2, Name: "painting"},
				{Id: 5, Name: "swimming"},
			},
		}
		search := "pa"

		mockRepo.EXPECT().GetHobbyBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		hobbyNames, err := s.GetHobbyBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, hobbyNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		search := "t"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Name: "coding"},
			{ID: 2, Name: "testing"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "coding"},
				{Id: 2, Name: "testing"},
			},
		}

		mockRepo.EXPECT().GetHobbyPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		hobbyNames, err := s.GetHobbyBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, hobbyNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		search := "dr"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetHobbyBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetHobbyBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetHobbyById(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		expectedHobbyName := "singing"

		var id int64 = 3

		mockRepo.EXPECT().GetHobbyByID(gomock.Any(), id).Return(expectedHobbyName, nil)

		s := service.NewService(mockRepo)
		hobbyName, err := s.GetHobbyById(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, hobbyName, &optionhubproto.GetByIdOut{Id: id, Value: expectedHobbyName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetHobbyByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetHobbyById(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddHobby(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		hobbyName := "boxing"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddHobby(gomock.Any(), hobbyName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddHobby(ctx, &optionhubproto.AddIn{Value: hobbyName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: hobbyName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		hobbyName := "boxing"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := service.NewMockDBRepo(ctrl)

		s := service.NewService(mockRepo)

		_, err := s.AddHobby(ctx, &optionhubproto.AddIn{Value: hobbyName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		hobbyName := "boxing"

		var expectedID int64

		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddHobby(gomock.Any(), hobbyName, uuid).Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddHobby(ctx, &optionhubproto.AddIn{Value: hobbyName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetSkillBySearchName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		expectedNames := []model.CategoryItem{
			{ID: 1, Name: "kafka"},
			{ID: 2, Name: "s3"},
			{ID: 5, Name: "swift"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "kafka"},
				{Id: 2, Name: "s3"},
				{Id: 5, Name: "swift"},
			},
		}
		search := "sw"

		mockRepo.EXPECT().GetSkillBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		skillNames, err := s.GetSkillBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, skillNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		search := "q"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Name: "sql"},
			{ID: 2, Name: "python"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "sql"},
				{Id: 2, Name: "python"},
			},
		}

		mockRepo.EXPECT().GetSkillPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		skillNames, err := s.GetSkillBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, skillNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		search := "go"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetSkillBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetSkillBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetSkillById(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		expectedSkillName := "kotlin"

		var id int64 = 3

		mockRepo.EXPECT().GetSkillByID(gomock.Any(), id).Return(expectedSkillName, nil)

		s := service.NewService(mockRepo)
		skillName, err := s.GetSkillById(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, skillName, &optionhubproto.GetByIdOut{Id: id, Value: expectedSkillName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetSkillByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetSkillById(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddSkill(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		skillName := "R"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddSkill(gomock.Any(), skillName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddSkill(ctx, &optionhubproto.AddIn{Value: skillName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: skillName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		skillName := "R"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := service.NewMockDBRepo(ctrl)

		s := service.NewService(mockRepo)

		_, err := s.AddSkill(ctx, &optionhubproto.AddIn{Value: skillName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		skillName := "R"

		var expectedID int64

		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddSkill(gomock.Any(), skillName, uuid).Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddSkill(ctx, &optionhubproto.AddIn{Value: skillName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetCityBySearchName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		expectedNames := []model.CategoryItem{
			{ID: 1, Name: "Moscow"},
			{ID: 2, Name: "New York"},
			{ID: 5, Name: "St. Petersburg"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "Moscow"},
				{Id: 2, Name: "New York"},
				{Id: 5, Name: "St. Petersburg"},
			},
		}
		search := "Mo"

		mockRepo.EXPECT().GetCityBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		cityNames, err := s.GetCityBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, cityNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		search := "v"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Name: "Voronezh"},
			{ID: 2, Name: "Vena"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Name: "Voronezh"},
				{Id: 2, Name: "Vena"},
			},
		}

		mockRepo.EXPECT().GetCityPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		cityNames, err := s.GetCityBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, cityNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		search := "ne"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetCityBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetCityBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetCityById(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		expectedCityName := "Almata"

		var id int64 = 3

		mockRepo.EXPECT().GetCityByID(gomock.Any(), id).Return(expectedCityName, nil)

		s := service.NewService(mockRepo)
		cityName, err := s.GetCityById(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, cityName, &optionhubproto.GetByIdOut{Id: id, Value: expectedCityName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetCityByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetCityById(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddCity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		cityName := "Dublin"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddCity(gomock.Any(), cityName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddCity(ctx, &optionhubproto.AddIn{Value: cityName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: cityName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		cityName := "Dublin"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := service.NewMockDBRepo(ctrl)

		s := service.NewService(mockRepo)

		_, err := s.AddCity(ctx, &optionhubproto.AddIn{Value: cityName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		cityName := "Dublin"

		var expectedID int64

		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddCity(gomock.Any(), cityName, uuid).Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddCity(ctx, &optionhubproto.AddIn{Value: cityName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}
