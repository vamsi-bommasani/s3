package terratests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// run runs a command in the repo root and returns stdout, stderr and error.
func run(ctx context.Context, dir string, name string, args ...string) (string, string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	fmt.Printf("COMMAND: %s %v\nSTDOUT:\n%s\nSTDERR:\n%s\n", name, args, outb.String(), errb.String())
	return outb.String(), errb.String(), err
}

func TestTerraformS3(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// repo root (one level up from terratests)
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("failed to determine repo root: %v", err)
	}

	// Initialize
	if out, errOut, err := run(ctx, repoRoot, "terraform", "init", "-input=false"); err != nil {
		t.Fatalf("terraform init failed: %v\nstdout: %s\nstderr: %s", err, out, errOut)
	}

	// Apply (use the test tfvars file)
	if out, errOut, err := run(ctx, repoRoot, "terraform", "apply", "-auto-approve", "-var-file=terratests/test.tfvars"); err != nil {
		t.Fatalf("terraform apply failed: %v\nstdout: %s\nstderr: %s", err, out, errOut)
	}

	// Ensure we always attempt to destroy at the end
	defer func() {
		if out, errOut, err := run(ctx, repoRoot, "terraform", "destroy", "-auto-approve", "-var-file=terratests/test.tfvars"); err != nil {
			t.Logf("terraform destroy failed: %v\nstdout: %s\nstderr: %s", err, out, errOut)
		}
	}()

	// Get bucket ID output
	out, errOut, err := run(ctx, repoRoot, "terraform", "output", "-json", "s3_bucket_id")
	if err != nil {
		t.Fatalf("terraform output failed: %v\nstdout: %s\nstderr: %s", err, out, errOut)
	}

	var bucketID string
	if err := json.Unmarshal([]byte(out), &bucketID); err != nil {
		t.Fatalf("failed to parse s3_bucket_id output: %v\noutput: %s", err, out)
	}

	if bucketID == "" {
		t.Fatalf("expected S3 bucket ID output, got empty string")
	}

	// Get bucket ARN output
	out, errOut, err = run(ctx, repoRoot, "terraform", "output", "-json", "s3_bucket_arn")
	if err != nil {
		t.Fatalf("terraform output failed: %v\nstdout: %s\nstderr: %s", err, out, errOut)
	}

	var bucketARN string
	if err := json.Unmarshal([]byte(out), &bucketARN); err != nil {
		t.Fatalf("failed to parse s3_bucket_arn output: %v\noutput: %s", err, out)
	}

	if bucketARN == "" {
		t.Fatalf("expected S3 bucket ARN output, got empty string")
	}

	// Get bucket region output
	out, errOut, err = run(ctx, repoRoot, "terraform", "output", "-json", "s3_bucket_region")
	if err != nil {
		t.Fatalf("terraform output failed: %v\nstdout: %s\nstderr: %s", err, out, errOut)
	}

	var bucketRegion string
	if err := json.Unmarshal([]byte(out), &bucketRegion); err != nil {
		t.Fatalf("failed to parse s3_bucket_region output: %v\noutput: %s", err, out)
	}

	if bucketRegion == "" {
		t.Fatalf("expected S3 bucket region output, got empty string")
	}

	t.Logf("Created S3 bucket: %s in region %s", bucketID, bucketRegion)
}
