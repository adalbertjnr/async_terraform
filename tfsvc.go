package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const Output = "output"

type TFService struct {
	execPath string
}

type installerFunc func() (string, error)

func NewTerraformService(installer installerFunc) (*TFService, error) {
	execPath, err := installer()
	if err != nil {
		return nil, err
	}
	return &TFService{
		execPath: execPath,
	}, nil
}

func terraformInstaller() (string, error) {
	installer := releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.7.0")),
	}
	return installer.Install(context.Background())
}

func (t *TFService) terraformTaskPlan(ctx context.Context, wd, execPath string) error {
	tf, err := tfexec.NewTerraform(wd, execPath)
	if err != nil {
		return err
	}

	tf.SetStdout(os.Stdout)
	if err := tf.Init(ctx, tfexec.Reconfigure(true)); err != nil {
		return err
	}

	_, err = tf.Plan(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *TFService) terraformTaskDestroy(ctx context.Context, wd, execPath string) error {
	tf, err := tfexec.NewTerraform(wd, execPath)
	if err != nil {
		return err
	}

	tf.SetStdout(os.Stdout)
	if err := tf.Init(ctx, tfexec.Reconfigure(true)); err != nil {
		return err
	}

	return tf.Destroy(ctx)
}

func (t *TFService) terraformTaskCreate(ctx context.Context, wd, execPath string) error {
	tf, err := tfexec.NewTerraform(wd, execPath)
	if err != nil {
		return err
	}

	tf.SetStdout(os.Stdout)
	if err := tf.Init(ctx, tfexec.Reconfigure(true)); err != nil {
		return err
	}

	plan, err := tf.Plan(ctx, tfexec.Out(Output))
	if err != nil {
		return err
	}

	if plan {
		err = tf.Apply(ctx, tfexec.DirOrPlan(Output))
		if err != nil {
			return err
		}
	} else {
		log.Printf("no changes detected in %s", wd)
	}
	return nil
}
