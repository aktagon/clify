package safety

import (
	"clify/internal/models"
	"regexp"
	"strings"
)

var (
	// Dangerous patterns that could cause system damage
	dangerousPatterns = []string{
		`rm\s+-rf\s+/`,
		`rm\s+-rf\s+\*`,
		`rm\s+.*\s+-rf`,
		`sudo\s+rm`,
		`chmod\s+777`,
		`chown\s+.*root`,
		`dd\s+if=.*of=/dev/`,
		`mkfs\s+`,
		`fdisk\s+`,
		`format\s+`,
		`del\s+/s\s+/q`,
		`rmdir\s+/s\s+/q`,
		`:\(\)\{.*\}`,        // Fork bomb
		`curl.*\|\s*sh`,
		`wget.*\|\s*sh`,
		`eval\s+.*\$\(`,
		`>\s*/dev/sd[a-z]`,
		`shutdown\s+`,
		`reboot\s+`,
		`halt\s+`,
		`poweroff\s+`,
		`killall\s+`,
		`pkill\s+.*-9`,
		`kill\s+.*-9.*1`,
		`systemctl\s+disable`,
		`systemctl\s+stop.*ssh`,
		`iptables\s+.*DROP`,
		`ufw\s+.*deny`,
		`passwd\s+root`,
		`useradd\s+.*sudo`,
		`usermod\s+.*sudo`,
		`crontab\s+.*root`,
		`/etc/passwd`,
		`/etc/shadow`,
		`/etc/sudoers`,
		`sudo\s+.*passwd`,
		`sudo\s+.*userdel`,
		`sudo\s+.*groupdel`,
	}

	// Warning patterns that modify system but are generally safe
	warningPatterns = []string{
		`sudo\s+`,
		`apt\s+install`,
		`apt\s+remove`,
		`apt\s+purge`,
		`yum\s+install`,
		`yum\s+remove`,
		`dnf\s+install`,
		`dnf\s+remove`,
		`pacman\s+-S`,
		`pacman\s+-R`,
		`brew\s+install`,
		`brew\s+uninstall`,
		`npm\s+install\s+-g`,
		`pip\s+install`,
		`pip\s+uninstall`,
		`cargo\s+install`,
		`go\s+install`,
		`chmod\s+`,
		`chown\s+`,
		`chgrp\s+`,
		`mkdir\s+`,
		`rmdir\s+`,
		`rm\s+`,
		`mv\s+`,
		`cp\s+`,
		`ln\s+`,
		`touch\s+`,
		`systemctl\s+start`,
		`systemctl\s+restart`,
		`systemctl\s+enable`,
		`service\s+`,
		`crontab\s+`,
		`export\s+`,
		`git\s+push`,
		`git\s+pull`,
		`git\s+clone`,
		`docker\s+run`,
		`docker\s+build`,
		`docker\s+pull`,
		`kubectl\s+apply`,
		`kubectl\s+delete`,
		`terraform\s+apply`,
		`terraform\s+destroy`,
		`>\s+`,
		`>>\s+`,
		`|\s+`,
		`&\s*$`,
		`nohup\s+`,
		`screen\s+`,
		`tmux\s+`,
	}

	// Safe patterns for read-only operations
	safePatterns = []string{
		`^ls\s+`,
		`^pwd\s*$`,
		`^whoami\s*$`,
		`^id\s*$`,
		`^date\s*$`,
		`^uptime\s*$`,
		`^uname\s+`,
		`^cat\s+`,
		`^less\s+`,
		`^more\s+`,
		`^head\s+`,
		`^tail\s+`,
		`^grep\s+`,
		`^find\s+`,
		`^locate\s+`,
		`^which\s+`,
		`^whereis\s+`,
		`^type\s+`,
		`^file\s+`,
		`^stat\s+`,
		`^du\s+`,
		`^df\s+`,
		`^free\s+`,
		`^ps\s+`,
		`^top\s+`,
		`^htop\s+`,
		`^jobs\s+`,
		`^history\s+`,
		`^env\s+`,
		`^printenv\s+`,
		`^echo\s+`,
		`^printf\s+`,
		`^wc\s+`,
		`^sort\s+`,
		`^uniq\s+`,
		`^cut\s+`,
		`^awk\s+`,
		`^sed\s+.*'s/`,
		`^tr\s+`,
		`^basename\s+`,
		`^dirname\s+`,
		`^realpath\s+`,
		`^readlink\s+`,
		`^git\s+status`,
		`^git\s+log`,
		`^git\s+diff`,
		`^git\s+show`,
		`^git\s+branch`,
		`^git\s+tag`,
		`^git\s+remote`,
		`^docker\s+ps`,
		`^docker\s+images`,
		`^docker\s+logs`,
		`^kubectl\s+get`,
		`^kubectl\s+describe`,
		`^kubectl\s+logs`,
		`^npm\s+list`,
		`^npm\s+outdated`,
		`^pip\s+list`,
		`^pip\s+show`,
		`^cargo\s+check`,
		`^cargo\s+test`,
		`^go\s+version`,
		`^go\s+env`,
		`^python\s+--version`,
		`^node\s+--version`,
		`^help\s+`,
		`^man\s+`,
		`^info\s+`,
		`--help\s*$`,
		`-h\s*$`,
	}
)

type Classifier struct {
	dangerousRegexes []*regexp.Regexp
	warningRegexes   []*regexp.Regexp
	safeRegexes      []*regexp.Regexp
}

func NewClassifier() *Classifier {
	c := &Classifier{}
	
	// Compile dangerous patterns
	for _, pattern := range dangerousPatterns {
		if regex, err := regexp.Compile("(?i)" + pattern); err == nil {
			c.dangerousRegexes = append(c.dangerousRegexes, regex)
		}
	}
	
	// Compile warning patterns
	for _, pattern := range warningPatterns {
		if regex, err := regexp.Compile("(?i)" + pattern); err == nil {
			c.warningRegexes = append(c.warningRegexes, regex)
		}
	}
	
	// Compile safe patterns
	for _, pattern := range safePatterns {
		if regex, err := regexp.Compile("(?i)" + pattern); err == nil {
			c.safeRegexes = append(c.safeRegexes, regex)
		}
	}
	
	return c
}

func (c *Classifier) ClassifyCommand(command string) models.SafetyLevel {
	// Normalize command for analysis
	normalized := strings.TrimSpace(command)
	
	// Check for dangerous patterns first
	for _, regex := range c.dangerousRegexes {
		if regex.MatchString(normalized) {
			return models.SafetyLevelDangerous
		}
	}
	
	// Check for safe patterns
	for _, regex := range c.safeRegexes {
		if regex.MatchString(normalized) {
			return models.SafetyLevelSafe
		}
	}
	
	// Check for warning patterns
	for _, regex := range c.warningRegexes {
		if regex.MatchString(normalized) {
			return models.SafetyLevelWarning
		}
	}
	
	// Default to warning for unknown commands
	return models.SafetyLevelWarning
}

func (c *Classifier) GetSafetyIcon(level models.SafetyLevel) string {
	switch level {
	case models.SafetyLevelSafe:
		return "🟢"
	case models.SafetyLevelWarning:
		return "🟡"
	case models.SafetyLevelDangerous:
		return "🔴"
	default:
		return "🟡"
	}
}

func (c *Classifier) GetSafetyMessage(level models.SafetyLevel) string {
	switch level {
	case models.SafetyLevelSafe:
		return "Safe to execute"
	case models.SafetyLevelWarning:
		return "Proceed with caution"
	case models.SafetyLevelDangerous:
		return "Dangerous command - use extreme caution"
	default:
		return "Unknown safety level"
	}
}