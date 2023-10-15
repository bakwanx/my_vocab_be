package controller

import (
	"encoding/json"
	"my_vocab/config"
	"my_vocab/dto/out"
	"my_vocab/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	filter = [2]string{"by_order", "by_search"}
)

func PostVocab(response http.ResponseWriter, request *http.Request) {
	var (
		result      out.Response
		vocabModels models.Vocab
	)
	response.Header().Set("Content-Type", "application/json")
	vocab := request.FormValue("vocab")
	idType := request.FormValue("id_type")
	idUser, _ := strconv.Atoi(request.FormValue("id_user"))
	variation := request.FormValue("variation")
	note := request.FormValue("note")
	timeNow := time.Now()

	// check user
	checkVocab := models.Vocab{}
	config.DB.Where("vocab = ?", vocab).First(&checkVocab)
	if checkVocab.Vocab != "" {
		response.WriteHeader(http.StatusConflict)
		result.Code = http.StatusConflict
		result.Status = "Failed"
		result.Message = "Vocab sudah terdaftar"
		json.NewEncoder(response).Encode(result)
		return
	}

	vocabModels = models.Vocab{
		IdUser:    idUser,
		Vocab:     vocab,
		IdType:    idType,
		Variation: variation,
		Note:      note,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := config.DB.Save(&vocabModels).Error

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status Internal Server Error"
		json.NewEncoder(response).Encode(result)
		return
	}

	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Data = vocabModels
	result.Message = "Berhasil menambahkan vocabulary"
	json.NewEncoder(response).Encode(result)
	return
}

func PatchVocab(response http.ResponseWriter, request *http.Request) {
	var (
		result      out.Response
		vocabModels models.Vocab
	)
	response.Header().Set("Content-Type", "application/json")
	idVocab := request.FormValue("id_vocab")
	vocab := request.FormValue("vocab")
	idType := request.FormValue("id_type")
	variation := request.FormValue("variation")
	note := request.FormValue("note")
	timeNow := time.Now()

	vocabModels = models.Vocab{
		Vocab:     vocab,
		IdType:    idType,
		Variation: variation,
		Note:      note,
		UpdatedAt: timeNow,
	}

	err := config.DB.Where("id_vocab = ?", idVocab).Updates(&vocabModels).Error

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status vocab tidak ditemukan"
		json.NewEncoder(response).Encode(result)
		return
	}

	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Data = vocabModels
	result.Message = "Berhasil update vocabulary"
	json.NewEncoder(response).Encode(result)
	return
}

func GetVocabularyByOrder(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var result out.Response

	idUser := mux.Vars(request)["id_user"]
	vocabModels := []models.Vocab{}

	// get by order alphabet
	config.DB.Model(models.Vocab{}).Where("id_user = ?", idUser).Order("vocab desc").Find(&vocabModels)
	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Data = vocabModels
	result.Message = "Berhasil mendapatkan vocabulary"
	json.NewEncoder(response).Encode(result)
	return

}

func GetVocabularyByDate(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var result out.Response

	idUser := mux.Vars(request)["id_user"]
	vocabModels := []models.Vocab{}

	// get by date
	config.DB.Model(models.Vocab{}).Where("id_user = ?", idUser).Order("created_at desc").Find(&vocabModels)
	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Data = vocabModels
	result.Message = "Berhasil mendapatkan vocabulary"
	json.NewEncoder(response).Encode(result)
	return

}

func GetVocabularyBySearch(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var result out.Response

	idUser := request.FormValue("id_user")
	keyword := request.FormValue("keyword")
	vocabModels := []models.Vocab{}

	// search by search keyword
	config.DB.Model(models.Vocab{}).Where("id_user = ? AND vocab LIKE ?%", idUser, keyword).Find(&vocabModels)
	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Data = vocabModels
	result.Message = "Berhasil mendapatkan vocabulary"
	json.NewEncoder(response).Encode(result)
	return

}

func PostTypeVocab(response http.ResponseWriter, request *http.Request) {
	var (
		result          out.Response
		typeVocabModels models.TypeVocab
	)
	response.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(request.Body).Decode(&typeVocabModels)
	timeNow := time.Now()

	// check user
	checkTypeVocab := models.TypeVocab{}
	config.DB.Where("type = ?", typeVocabModels.Type).First(&checkTypeVocab)
	if checkTypeVocab.Type != "" {
		response.WriteHeader(http.StatusConflict)
		result.Code = http.StatusConflict
		result.Status = "Failed"
		result.Message = "Vocab sudah terdaftar"
		json.NewEncoder(response).Encode(result)
		return
	}

	typeVocabModels = models.TypeVocab{
		Type:        typeVocabModels.Type,
		Description: typeVocabModels.Description,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	err = config.DB.Save(&typeVocabModels).Error

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status Internal Server Error"
		json.NewEncoder(response).Encode(result)
		return
	}

	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Data = typeVocabModels
	result.Message = "Berhasil menambahkan tipe vocabulary"
	json.NewEncoder(response).Encode(result)
	return
}

func GetType(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var result out.Response
	typeModels := []models.TypeVocab{}

	config.DB.Model(models.TypeVocab{}).Find(&typeModels)
	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Data = typeModels
	result.Message = "Berhasil mendapatkan tipe vocab"
	json.NewEncoder(response).Encode(result)
	return

}
