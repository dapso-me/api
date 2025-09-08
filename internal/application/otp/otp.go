package otp_usecase

import (
	"api/internal/common"
	"api/internal/domain/otp"
	"context"
	"fmt"
	"time"
)

type MailService interface {
	SendCode(to, templatePath, code string) error
}

type useCase struct {
	otpRepo     otp.Repository
	mailService MailService
}

func New(otpRepo otp.Repository, mailService MailService) otp.UseCase {
	return &useCase{
		otpRepo:     otpRepo,
		mailService: mailService,
	}
}

func (uc *useCase) CheckRateLimit(c context.Context, email, ip string) error {
	since := time.Now().UTC().Add(-time.Second * 60 * 5)

	// mail limit check
	otps, err := uc.otpRepo.FindByEmailSince(c, email, since)
	if err != nil {
		return fmt.Errorf("check rate limit: %w", err)
	}

	if len(otps) > 2 {
		return fmt.Errorf("check rate limit: %w", otp.ErrEmailRegistrationRateLimit)
	}

	// ip limit check
	otps, err = uc.otpRepo.FindByIPSince(c, email, since)
	if err != nil {
		return fmt.Errorf("check rate limit: %w", err)
	}

	if len(otps) > 4 {
		return fmt.Errorf("check rate limit: %w", common.ErrTooManyRequests)
	}

	return nil
}

func (uc *useCase) Send(c context.Context, dto *otp.SendInput) error {
	if err := uc.CheckRateLimit(c, dto.Email, dto.IP); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	otpEntity, err := otp.New(dto.Email, dto.Purpose, dto.IP)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	if err := uc.mailService.SendCode(dto.Email, string(dto.Purpose), otpEntity.Code); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (uc *useCase) VerifyAndRevoke(c context.Context, dto *otp.VerifyInput) error {
	otpEntity, err := uc.otpRepo.FindOneByCredentials(c, dto.Email, dto.Code)
	if err != nil {
		return fmt.Errorf("verify and revoke: %w", err)
	}

	if err := uc.otpRepo.Remove(c, otpEntity); err != nil {
		return fmt.Errorf("verify and revoke: %w", err)
	}

	return nil
}
