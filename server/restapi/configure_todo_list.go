// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"flightAPI/server/models"
	"flightAPI/server/repositories"
	"flightAPI/server/services"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"flightAPI/server/restapi/operations"
	"flightAPI/server/restapi/operations/todos"
)

//go:generate swagger generate server --target ..\..\flightAPI --name TodoList --spec ..\swagger.yaml --principal interface{}

func configureFlags(api *operations.TodoListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TodoListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	cfg, err := getServerConfig()

	ctx, mongoClient := connectMongoClient(cfg, err)

	var itemRepository repositories.ItemRepository = repositories.NewMongoItemRepository(mongoClient, cfg.Mongo.DatabaseName, "todo_list")
	var fileService services.FileService = services.NewLocalFileService()

	api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
		item, err := itemRepository.AddItem(*params.Body)
		if err != nil && err.Error() == "already_exist" {
			message := fmt.Sprintf("item already exist")
			payload := &models.Error{Code: http.StatusConflict, Message: &message}
			return todos.NewAddOneDefault(http.StatusConflict).WithPayload(payload)
		} else if err != nil {
			message := fmt.Sprintf(err.Error())
			payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
			return todos.NewAddOneDefault(http.StatusInternalServerError).WithPayload(payload)
		} else {
			return todos.NewAddOneCreated().WithPayload(item)
		}
	})

	api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
		err := itemRepository.DeleteItem(params.ID)
		if err != nil && err.Error() == "not_found" {
			message := fmt.Sprintf("item not found")
			payload := &models.Error{Code: http.StatusNotFound, Message: &message}
			return todos.NewDestroyOneDefault(http.StatusNotFound).WithPayload(payload)
		} else if err != nil {
			message := fmt.Sprintf(err.Error())
			payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
			return todos.NewDestroyOneDefault(http.StatusInternalServerError).WithPayload(payload)
		} else {
			return todos.NewDestroyOneNoContent()
		}
	})

	api.TodosGetOneHandler = todos.GetOneHandlerFunc(func(params todos.GetOneParams) middleware.Responder {
		item, err := itemRepository.FindItem(params.ID)
		if err != nil && err.Error() == "not_found" {
			message := fmt.Sprintf("item not found")
			payload := &models.Error{Code: http.StatusNotFound, Message: &message}
			return todos.NewGetOneDefault(http.StatusNotFound).WithPayload(payload)
		} else if err != nil {
			message := fmt.Sprintf(err.Error())
			payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
			return todos.NewGetOneDefault(http.StatusInternalServerError).WithPayload(payload)
		} else {
			return todos.NewGetOneOK().WithPayload(item)
		}
	})

	api.TodosFindTodosHandler = todos.FindTodosHandlerFunc(func(params todos.FindTodosParams) middleware.Responder {
		since, limit := int32(0), int32(10)
		if params.Since != nil {
			since = *params.Since
		}
		if params.Limit != nil {
			limit = *params.Limit
		}
		items, err := itemRepository.FindItems(since, limit)
		if err != nil && err.Error() == "not_found" {
			message := fmt.Sprintf("item not found")
			payload := &models.Error{Code: http.StatusNotFound, Message: &message}
			return todos.NewFindTodosDefault(http.StatusNotFound).WithPayload(payload)
		} else if err != nil {
			message := fmt.Sprintf(err.Error())
			payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
			return todos.NewFindTodosDefault(http.StatusInternalServerError).WithPayload(payload)
		} else {
			return todos.NewFindTodosOK().WithPayload(items)
		}
	})

	api.TodosUpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams) middleware.Responder {
		item := *params.Body
		item.ID = params.ID
		itemResult, err := itemRepository.UpdateItem(item)
		if err != nil && err.Error() == "not_found" {
			message := fmt.Sprintf("Item Not Found")
			payload := &models.Error{Code: http.StatusNotFound, Message: &message}
			return todos.NewUpdateOneDefault(http.StatusNotFound).WithPayload(payload)
		} else if err != nil {
			message := fmt.Sprintf(err.Error())
			payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
			return todos.NewUpdateOneDefault(http.StatusInternalServerError).WithPayload(payload)
		} else {
			return todos.NewUpdateOneOK().WithPayload(itemResult)
		}
	})

	api.TodosUploadFileHandler = todos.UploadFileHandlerFunc(func(params todos.UploadFileParams) middleware.Responder {
		file := params.File
		defer file.Close()
		fileMeta, err := fileService.UploadFile(file)
		if err != nil {
			message := fmt.Sprintf("Failed To Upload the file: %s", err.Error())
			payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
			return todos.NewUploadFileDefault(http.StatusInternalServerError).WithPayload(payload)
		}
		return todos.NewUploadFileCreated().WithPayload(fileMeta)
	})

	api.TodosDownloadFileHandler = todos.DownloadFileHandlerFunc(func(params todos.DownloadFileParams) middleware.Responder {
		uuid := params.UUID
		file, err := fileService.DownloadFile(uuid)
		if err != nil {
			if err.Error() == "not_found" {
				message := fmt.Sprintf("File Not Found")
				payload := &models.Error{Code: http.StatusNotFound, Message: &message}
				return todos.NewDownloadFileDefault(http.StatusNotFound).WithPayload(payload)
			} else {
				message := fmt.Sprintf("Cannot delete file: %s", err)
				payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
				return todos.NewDownloadFileDefault(http.StatusInternalServerError).WithPayload(payload)
			}
		}
		response := todos.NewDownloadFileOK()
		response.Payload = file
		return response
	})

	api.TodosDeleteFileHandler = todos.DeleteFileHandlerFunc(func(params todos.DeleteFileParams) middleware.Responder {
		uuid := params.UUID
		if err := fileService.DeleteFile(uuid); err != nil {
			if err.Error() == "not_found" {
				message := fmt.Sprintf("File Not Found")
				payload := &models.Error{Code: http.StatusNotFound, Message: &message}
				return todos.NewDeleteFileDefault(http.StatusNotFound).WithPayload(payload)
			} else {
				message := fmt.Sprintf("Cannot delete file: %s", err)
				payload := &models.Error{Code: http.StatusInternalServerError, Message: &message}
				return todos.NewDeleteFileDefault(http.StatusInternalServerError).WithPayload(payload)
			}
		}
		return todos.NewDeleteFileNoContent()
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
		api.Logger("Disconnected Mongo")
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func connectMongoClient(cfg ServerConfig, err error) (context.Context, *mongo.Client) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s",
		cfg.Mongo.User, cfg.Mongo.Password, cfg.Mongo.ConnectionString, cfg.Mongo.DatabaseName)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return ctx, client
}

func getServerConfig() (ServerConfig, error) {
	// config
	f, err := os.Open("C:\\Users\\vanme\\go\\src\\flightAPI\\config\\server-config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg ServerConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
