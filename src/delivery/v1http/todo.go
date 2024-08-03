package v1http

import (
	"encoding/json"
	"errors"
	"fmt"
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
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
		case http.MethodGet:
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, constants.ErrMsgEmptyId)
			return
		}

		switch r.Method {
		case http.MethodGet:
			var record model.TaskPresenter
			var mRecord []byte

			fn.Do(func() (err error) {
				record, err = h.Useacse.GetTask(id)
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
				w.WriteHeader(msgError.StatusCode)
				fmt.Fprintf(w, msgError.Message)
				return
			} else if err != nil {
				delivery.HandleUnhandledError(w, err, logger)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderApplicationJSON)
			fmt.Fprint(w, string(mRecord))
			return
		case http.MethodPut:
		case http.MethodDelete:
			fn := new(common.LeastError)

			isHardDeleteStr := r.Header.Get("hard-delete")
			isHardDelete, _ := strconv.ParseBool(isHardDeleteStr)

			fn.Do(func() (err error) {
				err = h.Useacse.DeleteTask(id, isHardDelete)
				return
			})

			err := fn.Err()
			asMsgError := errors.As(err, &msgError)

			if asMsgError {
				w.WriteHeader(msgError.StatusCode)
				fmt.Fprintf(w, msgError.Message)
				return
			} else if err != nil {
				delivery.HandleUnhandledError(w, err, logger)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
			fmt.Fprint(w, constants.MsgSuccess)
			return
		default:
			delivery.HandleUnknownHttpMethod(w)
			return
		}
	}
}
