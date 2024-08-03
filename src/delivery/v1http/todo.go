package v1http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sagala-todo/pkg/common"
	"sagala-todo/pkg/constants"
	customerror "sagala-todo/pkg/custom-error"
	"sagala-todo/src/delivery"
	"sagala-todo/src/model"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (h *V1Handler) RootHandler() http.HandlerFunc {
	const handlerName = "RootHandler"
	logger := logger.WithFields(logrus.Fields{
		constants.LFKHandlerName: handlerName,
	})
	return func(w http.ResponseWriter, r *http.Request) {
		var msgError *customerror.HttpError
		fn := new(common.LeastError)

		switch r.Method {
		case http.MethodPost:
			var (
				taskId  string
				payload model.TaskDTO
			)

			fn.Do(func() (err error) {
				bt, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Error(err)
					err = &customerror.HttpError{
						Message:    fmt.Sprintf("failed read body request: %s", constants.ErrCodeFailedReadBodyRequest),
						StatusCode: http.StatusInternalServerError,
						Err:        err,
					}
					return
				}

				err = json.Unmarshal(bt, &payload)
				if err != nil {
					logger.Error(err)
					err = &customerror.HttpError{
						Message:    fmt.Sprintf("invalid payload: %s", constants.ErrInvalidPayload),
						StatusCode: http.StatusBadRequest,
						Err:        err,
					}
					return
				}
				return
			})

			fn.Do(func() (err error) {
				taskId, err = h.Usecase.PostTask(payload)
				return
			})

			err := fn.Err()
			asMsgError := errors.As(err, &msgError)

			if asMsgError {
				w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
				w.WriteHeader(msgError.StatusCode)
				fmt.Fprintf(w, msgError.Message)
				return
			} else if err != nil {
				delivery.HandleUnhandledError(w, err, logger)
				return
			}

			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, taskId)
			return
		case http.MethodGet:
			paramsLookup := []string{"limit", "offset", "search", "status"}
			queryParams := r.URL.Query()

			var (
				limit, offset   *int
				search, status  *string
				tasksPresenter  []model.TaskPresenter
				totalTasks      int
				responseContent []byte
			)

			fn.Do(func() (err error) {
				if ok := queryParams.Has(paramsLookup[0]); ok {
					i, _err := strconv.Atoi(queryParams.Get(paramsLookup[0]))
					if _err != nil {
						err = &customerror.HttpError{
							Message:    "limit should be a number",
							StatusCode: http.StatusBadRequest,
							Err:        _err,
						}
						return
					}
					limit = &i
				} else if !ok {
					err = &customerror.HttpError{
						Message:    "limit could not be empty",
						StatusCode: http.StatusBadRequest,
					}
					return
				}
				return
			})

			fn.Do(func() (err error) {
				if ok := queryParams.Has(paramsLookup[1]); ok {
					i, _err := strconv.Atoi(queryParams.Get(paramsLookup[1]))
					if _err != nil {
						err = &customerror.HttpError{
							Message:    "offset should be a number",
							StatusCode: http.StatusBadRequest,
							Err:        _err,
						}
						return
					}
					offset = &i
				} else if !ok {
					err = &customerror.HttpError{
						Message:    "offset could not be empty",
						StatusCode: http.StatusBadRequest,
					}
					return
				}
				return
			})

			fn.Do(func() (err error) {
				if ok := queryParams.Has(paramsLookup[2]); ok {
					search = new(string)
					*search = queryParams.Get(paramsLookup[2])
				}
				return
			})

			fn.Do(func() (err error) {
				if ok := queryParams.Has(paramsLookup[3]); ok {
					status = new(string)
					*status = queryParams.Get(paramsLookup[3])
					logger.Info(*status)
				}
				return
			})

			fn.Do(func() (err error) {
				tasksPresenter, totalTasks, err = h.Usecase.GetTasks(limit, offset, search, status)
				return
			})

			fn.Do(func() (err error) {
				bt, err := json.Marshal(tasksPresenter)
				if err != nil {
					err = &customerror.HttpError{
						Message:    constants.ErrMsgFailedSerializeRecord,
						StatusCode: http.StatusInternalServerError,
						Err:        err,
					}
					return
				}
				resp := json.RawMessage(fmt.Sprintf("{\"message\":\"success\", \"data\":%s, \"total_data\":%d}", string(bt), totalTasks))
				responseContent, err = resp.MarshalJSON()
				if err != nil {
					responseContent = nil
					err = &customerror.HttpError{
						Message:    constants.ErrMsgFailedSerializeRecord,
						StatusCode: http.StatusInternalServerError,
						Err:        err,
					}
					return
				}
				return
			})

			err := fn.Err()
			asMsgError := errors.As(err, &msgError)

			if asMsgError {
				w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
				w.WriteHeader(msgError.StatusCode)
				fmt.Fprintf(w, msgError.Message)
				return
			} else if err != nil {
				delivery.HandleUnhandledError(w, err, logger)
				return
			}

			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderApplicationJson)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(responseContent))
			return
		default:
			delivery.HandleUnknownHttpMethod(w)
			return
		}
	}
}

func (h *V1Handler) TaskDetail() http.HandlerFunc {
	const handlerName = "TaskDetail"
	logger := logger.WithFields(logrus.Fields{
		constants.LFKHandlerName: handlerName,
	})
	return func(w http.ResponseWriter, r *http.Request) {
		var msgError *customerror.HttpError
		fn := new(common.LeastError)

		id := r.PathValue("id")
		if id == "" {
			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, constants.ErrMsgEmptyId)
			return
		}

		switch r.Method {
		case http.MethodGet:
			var record model.TaskPresenter
			var mRecord []byte

			fn.Do(func() (err error) {
				record, err = h.Usecase.GetTask(id)
				return
			})

			fn.Do(func() (err error) {
				mRecord, err = json.Marshal(record)
				if err != nil {
					logger.Error(err)
					err = &customerror.HttpError{
						Message:    constants.ErrMsgFailedSerializeRecord,
						Err:        err,
						StatusCode: http.StatusInternalServerError,
					}
				}
				return
			})

			err := fn.Err()
			asMsgError := errors.As(err, &msgError)

			if asMsgError {
				w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
				w.WriteHeader(msgError.StatusCode)
				fmt.Fprintf(w, msgError.Message)
				return
			} else if err != nil {
				delivery.HandleUnhandledError(w, err, logger)
				return
			}

			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderApplicationJSON)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(mRecord))
			return
		case http.MethodPut:
			var payload model.TaskDTO

			fn.Do(func() (err error) {
				bt, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Error(err)
					err = &customerror.HttpError{
						Message:    fmt.Sprintf("failed read body request: %s", constants.ErrCodeFailedReadBodyRequest),
						StatusCode: http.StatusInternalServerError,
						Err:        err,
					}
					return
				}

				err = json.Unmarshal(bt, &payload)
				if err != nil {
					logger.Error(err)
					err = &customerror.HttpError{
						Message:    fmt.Sprintf("invalid payload: %s", constants.ErrInvalidPayload),
						StatusCode: http.StatusBadRequest,
						Err:        err,
					}
					return
				}
				return
			})

			fn.Do(func() (err error) {
				err = h.Usecase.UpdateTask(id, payload)
				return
			})

			err := fn.Err()
			asMsgError := errors.As(err, &msgError)

			if asMsgError {
				w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
				w.WriteHeader(msgError.StatusCode)
				fmt.Fprintf(w, msgError.Message)
				return
			} else if err != nil {
				delivery.HandleUnhandledError(w, err, logger)
				return
			}

			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, constants.MsgSuccess)
			return
		case http.MethodDelete:
			isHardDeleteStr := r.Header.Get("hard-delete")
			isHardDelete, _ := strconv.ParseBool(isHardDeleteStr)

			fn.Do(func() (err error) {
				err = h.Usecase.DeleteTask(id, isHardDelete)
				return
			})

			err := fn.Err()
			asMsgError := errors.As(err, &msgError)

			if asMsgError {
				w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
				w.WriteHeader(msgError.StatusCode)
				fmt.Fprintf(w, msgError.Message)
				return
			} else if err != nil {
				delivery.HandleUnhandledError(w, err, logger)
				return
			}

			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, constants.MsgSuccess)
			return
		default:
			delivery.HandleUnknownHttpMethod(w)
			return
		}
	}
}
