package base

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/command"
	"github.com/CedricThomas/console/internal/service/keystore"
	"github.com/CedricThomas/console/internal/usecase/boot"
)

const bootOSKey = "boot:os:target"

type bootUsecase struct {
	keystore         keystore.Keystore
	cmdExec          command.CommandExecutor
	bootOSTTLSeconds time.Duration
}

func New(ks keystore.Keystore, ce command.CommandExecutor, bootOSTTLSeconds int) boot.Boot {
	return &bootUsecase{keystore: ks, cmdExec: ce, bootOSTTLSeconds: time.Duration(bootOSTTLSeconds) * time.Second}
}

func (b *bootUsecase) StoreBootOS(ctx context.Context, osName domain.OSName) error {
	if osName == "" {
		return errors.New("os name is required")
	}
	return b.keystore.SetWithTTL(ctx, bootOSKey, string(osName), b.bootOSTTLSeconds)
}

func (b *bootUsecase) GetBootOS(ctx context.Context) (domain.OSName, error) {
	val, err := b.keystore.Get(ctx, bootOSKey)
	if err != nil {
		return "", fmt.Errorf("get boot OS: %w", err)
	}
	if val == "" {
		return "", errors.New("no target OS stored")
	}
	return domain.OSName(val), nil
}

func (b *bootUsecase) RebootToOS(ctx context.Context) error {
	osName, err := b.GetBootOS(ctx)
	if err != nil {
		return err
	}
	entryName, err := b.matchGrubEntryToOS(osName)
	if err != nil {
		return fmt.Errorf("match GRUB entry: %w", err)
	}
	if err := b.cmdExec.SetGrubReboot(ctx, entryName); err != nil {
		var unsupported *domain.ErrUnsupportedOS
		if errors.As(err, &unsupported) {
			return &domain.ErrUnsupportedOS{}
		}
		return fmt.Errorf("set grub reboot: %w", err)
	}
	return b.cmdExec.Reboot(ctx)
}

func (b *bootUsecase) ListGrubEntries(ctx context.Context) ([]domain.BootEntry, error) {
	return b.cmdExec.ListGrubEntries(ctx)
}

func (b *bootUsecase) MatchAndRebootToOS(ctx context.Context, osName domain.OSName) error {
	entryName, err := b.matchGrubEntryToOS(osName)
	if err != nil {
		return err
	}
	if err := b.cmdExec.SetGrubReboot(ctx, entryName); err != nil {
		var unsupported *domain.ErrUnsupportedOS
		if errors.As(err, &unsupported) {
			return &domain.ErrUnsupportedOS{}
		}
		return fmt.Errorf("set grub reboot: %w", err)
	}
	return b.cmdExec.Reboot(ctx)
}

func (b *bootUsecase) matchGrubEntryToOS(osName domain.OSName) (string, error) {
	entries, err := b.cmdExec.ListGrubEntries(context.Background())
	if err != nil {
		return "", fmt.Errorf("list GRUB entries: %w", err)
	}
	return domain.MatchGrubEntryToOS(entries, osName)
}
