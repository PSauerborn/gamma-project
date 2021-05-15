package filestore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
)

// API handler used to serve health check handler
func HealthCheckHandler(ctx *gin.Context) {
	log.Info("received request for health check handler")
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Service running"})
}

// API handler user to retrieve a given file
func GetFileHandler(ctx *gin.Context) {
	log.Info("received request to retrieve file")
	fileId, err := uuid.Parse(ctx.Param("fileId"))
	if err != nil {
		log.Error(fmt.Errorf("received invalid file ID %s", ctx.Param("fileId")))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid file ID"})
		return
	}
	// retrieve file metadata from persistence layer
	file, err := persistence.GetFileMetadata(fileId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve file metadata: %+v", err))
		switch err {
		case ErrFileNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find file"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	// open file with given file path
	contents, err := persistence.GetFileContents(file)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve file contents: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	// generate reader from response body and attach to response
	responseBody := bytes.NewReader(contents)
	ctx.DataFromReader(http.StatusOK, int64(len(contents)), "", responseBody, nil)
}

// API handler user to retrieve all file metadata
// from the persistence layer
func ListFilesHandler(ctx *gin.Context) {
	log.Info("received request to retrieve metadata for all files")
	// retrieve metadata for all files from persistence layer
	files, err := persistence.ListFiles()
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve file(s): %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"files": files, "count": len(files)})
}

// API handler user to retrieve a given file
func GetFileMetadataHandler(ctx *gin.Context) {
	log.Info("received request to retrieve file")
	fileId, err := uuid.Parse(ctx.Param("fileId"))
	if err != nil {
		log.Error(fmt.Errorf("received invalid file ID %s", ctx.Param("fileId")))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid file ID"})
		return
	}
	// retrieve file metadata from persistence layer
	file, err := persistence.GetFileMetadata(fileId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve file metadata: %+v", err))
		switch err {
		case ErrFileNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find file"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"metadata": file})
}

// API handler used to create a new file
func CreateFileHandler(ctx *gin.Context) {
	log.Info("received request to create file")
	var request struct {
		Meta     map[string]interface{} `json:"meta" binding:"required"`
		FileName string                 `json:"file_name" binding:"required"`
		Content  string                 `json:"content" binding:"required"`
	}
	// extract request body from JSON content
	if err := ctx.ShouldBind(&request); err != nil {
		log.Error(fmt.Errorf("received invalid request body: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid request body"})
		return
	}
	// create new file instance
	body, err := utils.Base64ToBytes(request.Content)
	if err != nil {
		log.Error(fmt.Errorf("unable to decode base64 file string: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid request body"})
		return
	}
	// create new file instance via persistence interface
	fileId, err := persistence.CreateFile(body, request.FileName,
		request.Meta)
	if err != nil {
		log.Error(fmt.Errorf("unable to create new file instance: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
		"file_id": fileId})
}

// API handler used to modify an existing file
func PutFileHandler(ctx *gin.Context) {
	log.Info("received request to modify file")
	fileId, err := uuid.Parse(ctx.Param("fileId"))
	if err != nil {
		log.Error(fmt.Errorf("received invalid file ID %s", ctx.Param("fileId")))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid file ID"})
		return
	}
	// retrieve file metadata from persistence layer
	meta, err := persistence.GetFileMetadata(fileId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve file metadata: %+v", err))
		switch err {
		case ErrFileNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find file"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	// extract request body and read
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(fmt.Errorf("unable to extract request body: %+v", err))
		ctx.JSON(http.StatusBadRequest,
			gin.H{"http_code": http.StatusBadRequest, "message": "Invalid request body"})
		return
	}
	// mofidy file via persistence layer
	if err := persistence.ModifyFile(meta, body); err != nil {
		log.Error(fmt.Errorf("unable to modify file: %+v", err))
		switch err {
		case ErrFeatureNotSupported:
			ctx.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"http_code": http.StatusNotImplemented,
				"message": "Internal server error"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully modified file"})
}

// API handler used to delete a given file
func DeleteFileHandler(ctx *gin.Context) {
	log.Info("received request to delete file")
	fileId, err := uuid.Parse(ctx.Param("fileId"))
	if err != nil {
		log.Error(fmt.Errorf("received invalid file ID %s", ctx.Param("fileId")))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid file ID"})
		return
	}
	// retrieve file metadata from persistence layer
	meta, err := persistence.GetFileMetadata(fileId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve file metadata: %+v", err))
		switch err {
		case ErrFileNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find file"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	// delete file from persistence layer
	if err := persistence.DeleteFile(meta); err != nil {
		log.Error(fmt.Errorf("unable to delete file: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully deleted file"})
}

func ArchiveFileHandler(ctx *gin.Context) {
	log.Info("received request to archive file")
	fileId, err := uuid.Parse(ctx.Param("fileId"))
	if err != nil {
		log.Error(fmt.Errorf("received invalid file ID %s", ctx.Param("fileId")))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid file ID"})
		return
	}
	// retrieve file metadata from persistence layer
	meta, err := persistence.GetFileMetadata(fileId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve file metadata: %+v", err))
		switch err {
		case ErrFileNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find file"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}

	if err := persistence.ArchiveFile(meta); err != nil {
		log.Error(fmt.Errorf("unable to archive file: %+v", err))
		switch err {
		case ErrFeatureNotSupported:
			ctx.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"http_code": http.StatusNotImplemented,
				"message": "Archiving not supported"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully archived file"})
}

func SearchFilesHandler(ctx *gin.Context) {
	log.Info("received request for search")
	var request struct {
		SearchTerms map[string]interface{} `json:"search_terms"`
	}
	if err := ctx.ShouldBind(&request); err != nil {
		log.Error(fmt.Errorf("received invalid request body: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid search request"})
		return
	}
	// search files by metadata
	results, err := persistence.SearchFilesByMetadata(request.SearchTerms)
	if err != nil {
		log.Error(fmt.Errorf("unable to search files: %+v", err))
		switch err {
		case ErrFeatureNotSupported:
			ctx.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"http_code": http.StatusNotImplemented,
				"message": "Searching not supported"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"results": results})
}
