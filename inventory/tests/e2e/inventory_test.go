//go:build integration

package integration

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ploskirev/go-rocket/platform/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	inventoryV1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		// ufoClient ufoV1.UFOServiceClient
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		logger.Info(ctx, fmt.Sprintf("Adress: %s", env.App.Address()))

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	// Describe("Create", func() {
	// 	It("должен успешно создавать новую деталь", func() {
	// 		info := env.GetTestPartInfo()

	// 		resp, err := inventoryClient..Create(ctx, &ufoV1.CreateRequest{
	// 			Info: info,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(resp.GetUuid()).ToNot(BeEmpty())
	// 		Expect(resp.GetUuid()).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))
	// 	})
	// })

	Describe("Get", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовую деталь
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
		})

		It("Mock TEST", func() {
			Expect(true).To(BeTrue())
		})
		It("Mock TEST 2", func() {
			Expect(true).To(BeTrue())
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred(), "ожидали успешное получение детали")
			Expect(resp.GetUuid()).To(Equal(partUUID))
			Expect(resp.GetName()).ToNot(BeNil())
			Expect(resp.GetPrice()).ToNot(BeNil())
		})
	})

	// Describe("Update", func() {
	// 	var sightingUUID string

	// 	BeforeEach(func() {
	// 		// Вставляем тестовое наблюдение
	// 		var err error
	// 		sightingUUID, err = env.InsertTestSighting(ctx)
	// 		Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестового наблюдения в MongoDB")
	// 	})

	// 	It("должен успешно обновлять наблюдение", func() {
	// 		updateInfo := env.GetUpdatedSightingInfo()

	// 		_, err := ufoClient.Update(ctx, &ufoV1.UpdateRequest{
	// 			Uuid:       sightingUUID,
	// 			UpdateInfo: updateInfo,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())

	// 		// Проверяем, что наблюдение действительно обновилось
	// 		resp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
	// 			Uuid: sightingUUID,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(resp.GetSighting().GetInfo().Location).To(Equal(updateInfo.Location.GetValue()))
	// 		Expect(resp.GetSighting().GetInfo().Description).To(Equal(updateInfo.Description.GetValue()))
	// 		Expect(resp.GetSighting().GetInfo().Color.GetValue()).To(Equal(updateInfo.Color.GetValue()))
	// 		Expect(resp.GetSighting().GetInfo().DurationSeconds.GetValue()).To(Equal(updateInfo.DurationSeconds.GetValue()))
	// 		Expect(resp.GetSighting().GetUpdatedAt()).ToNot(BeNil())
	// 	})
	// })

	// Describe("Delete", func() {
	// 	var sightingUUID string

	// 	BeforeEach(func() {
	// 		// Вставляем тестовое наблюдение
	// 		var err error
	// 		sightingUUID, err = env.InsertTestSighting(ctx)
	// 		Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестового наблюдения в MongoDB")
	// 	})

	// 	It("должен успешно выполнять мягкое удаление наблюдения", func() {
	// 		_, err := ufoClient.Delete(ctx, &ufoV1.DeleteRequest{
	// 			Uuid: sightingUUID,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())

	// 		// Проверяем, что наблюдение помечено как удаленное
	// 		resp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
	// 			Uuid: sightingUUID,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(resp.GetSighting().GetDeletedAt()).ToNot(BeNil())
	// 	})
	// })

	// Describe("Полный жизненный цикл", func() {
	// 	It("должен поддерживать полный CRUD цикл", func() {
	// 		// 1. Создаем наблюдение
	// 		info := env.GetTestSightingInfo()
	// 		createResp, err := ufoClient.Create(ctx, &ufoV1.CreateRequest{
	// 			Info: info,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(createResp.GetUuid()).ToNot(BeEmpty())
	// 		uuid := createResp.GetUuid()

	// 		// 2. Получаем созданное наблюдение
	// 		getResp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
	// 			Uuid: uuid,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(getResp.GetSighting().Uuid).To(Equal(uuid))
	// 		Expect(getResp.GetSighting().GetInfo().Location).To(Equal(info.Location))
	// 		Expect(getResp.GetSighting().GetInfo().Description).To(Equal(info.Description))

	// 		// 3. Обновляем наблюдение
	// 		updateInfo := env.GetUpdatedSightingInfo()
	// 		_, err = ufoClient.Update(ctx, &ufoV1.UpdateRequest{
	// 			Uuid:       uuid,
	// 			UpdateInfo: updateInfo,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())

	// 		// 4. Проверяем обновление
	// 		getUpdatedResp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
	// 			Uuid: uuid,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(getUpdatedResp.GetSighting().GetInfo().Location).To(Equal(updateInfo.Location.GetValue()))
	// 		Expect(getUpdatedResp.GetSighting().GetInfo().Description).To(Equal(updateInfo.Description.GetValue()))

	// 		// 5. Удаляем наблюдение
	// 		_, err = ufoClient.Delete(ctx, &ufoV1.DeleteRequest{
	// 			Uuid: uuid,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())

	// 		// 6. Проверяем, что наблюдение помечено как удаленное
	// 		getDeletedResp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
	// 			Uuid: uuid,
	// 		})

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(getDeletedResp.GetSighting().GetDeletedAt()).ToNot(BeNil())
	// 	})
	// })
})
