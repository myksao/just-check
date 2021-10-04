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

func (justChkUC justCheckUseCase) Check(ctx context.Context,file multipart.File,contentType string)  (string,error){

	// open and decode image file
	//file, _ := os.Open("qrcode.jpg")
	img, _, imgErr := image.Decode(file)
	if imgErr != nil {
		justChkUC.logger.Sugar().Error(imgErr)
		return "",errors.Wrap(imgErr,"Couldn't decode the image")
	}
	// prepare BinaryBitmap
	bmp, bmpErr := gozxing.NewBinaryBitmapFromImage(img)
	if bmpErr != nil {
		justChkUC.logger.Sugar().Error(bmpErr)
		return "",errors.Wrap(bmpErr,"Couldn't convert the image")
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, resultErr := qrReader.Decode(bmp, nil)
	if resultErr != nil {
		justChkUC.logger.Sugar().Error(resultErr)
		return "", errors.Wrap(resultErr,"Couldn't read the data embedded in the code")
	}

	if len(contentType)==0{
		validatorErr := justChkUC.validator.Var(result.GetText(), fmt.Sprintf("%s,%s", "required","ascii" ))
		if validatorErr != nil {
			justChkUC.logger.Sugar().Error(validatorErr)
			return "", errors.Wrap(validatorErr, "Couldn't validate the content type of the QRCode")
		}
	}else {
		validatorErr := justChkUC.validator.Var(result.GetText(), fmt.Sprintf("%s,%s", "required", strings.ToLower(contentType)))
		if validatorErr != nil {
			justChkUC.logger.Sugar().Error(validatorErr)
			return "", errors.Wrap(validatorErr, "Couldn't validate the content type of the QRCode")
		}
	}

	return result.GetText(),nil
}