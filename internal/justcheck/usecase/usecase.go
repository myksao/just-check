package usecase

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime/multipart"
	"strings"
)
type justCheckUseCase struct {
	logger *zap.Logger
	validator *validator.Validate
}


func NewJustCheckUseCase(logger *zap.Logger,validator *validator.Validate) *justCheckUseCase{
	return &justCheckUseCase{
		logger: logger,
		validator: validator,
	}
}

func (justChkUC justCheckUseCase) Check(ctx context.Context,file multipart.File,contentType string)  (result string, err error){

	defer func()  {
		if pErr := recover(); pErr != nil {
			log.Println("panic occurred:", pErr)
			switch x := pErr.(type) {
			case string:
				//TODO: Add a custom error message
				err = errors.New(x)

			case error:
				err = x

			default:
				err = errors.New("Unknown panic")

			}

		}
	}()

	// open and decode image file
	//file, _ := os.Open("qrcode.jpg")
	img, _, imgErr := image.Decode(file)
	if imgErr != nil {
		justChkUC.logger.Sugar().Error(imgErr)
		err = errors.Wrap(imgErr,"Couldn't decode the image")
		return result,err
	}
	// prepare BinaryBitmap
	bmp, bmpErr := gozxing.NewBinaryBitmapFromImage(img)
	if bmpErr != nil {
		justChkUC.logger.Sugar().Error(bmpErr)
		err = errors.Wrap(bmpErr,"Couldn't convert the image")
		return result,err
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	resultQrReader, resultErr := qrReader.Decode(bmp, nil)
	if resultErr != nil {
		justChkUC.logger.Sugar().Error(resultErr)

		err = errors.Wrap(resultErr,"Couldn't read the data embedded in the code")
		return result,err
	}
	result = resultQrReader.GetText()

	if len(contentType) > 0 {
		validatorErr := justChkUC.validator.Var(result, fmt.Sprintf("%s,%s", "required", strings.ToLower(contentType)))
		if validatorErr != nil {
			justChkUC.logger.Sugar().Error(validatorErr)
			err = errors.Wrap(validatorErr, "Couldn't validate the content type of the QRCode")
			return result,err
		}

	}else {

		validatorErr := justChkUC.validator.Var(result, fmt.Sprintf("%s,%s", "required","ascii" ))
		if validatorErr != nil {
			justChkUC.logger.Sugar().Error(validatorErr)
			err = errors.Wrap(validatorErr, "Couldn't validate the content type of the QRCode")
			return result,err
		}
	}

	return result,err
}

