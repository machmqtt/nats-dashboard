package store

import (
	"testing"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestCreateAndAuthenticate(t *testing.T) {
	s := testStore(t)

	u, err := s.CreateUser("admin", "secret123", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if u.Username != "admin" {
		t.Errorf("username = %q, want admin", u.Username)
	}
	if u.Role != RoleAdmin {
		t.Errorf("role = %q, want %q", u.Role, RoleAdmin)
	}

	authed, err := s.Authenticate("admin", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	if authed.ID != u.ID {
		t.Errorf("authenticated user ID mismatch")
	}
	if authed.Role != RoleAdmin {
		t.Errorf("authenticated role = %q, want %q", authed.Role, RoleAdmin)
	}
	if authed.LastLogin == nil {
		t.Error("expected last_login to be set after successful auth")
	}
}

func TestCreateUserInvalidRole(t *testing.T) {
	s := testStore(t)
	_, err := s.CreateUser("user", "pass", "superadmin")
	if err == nil {
		t.Fatal("expected error for invalid role")
	}
}

func TestAuthenticateBadPassword(t *testing.T) {
	s := testStore(t)
	s.CreateUser("admin", "secret123", RoleAdmin)

	_, err := s.Authenticate("admin", "wrong")
	if err == nil {
		t.Fatal("expected error for bad password")
	}

	// Check failed attempt was recorded.
	u, _ := s.GetUser(1)
	if u.FailedAttempts != 1 {
		t.Errorf("failed_attempts = %d, want 1", u.FailedAttempts)
	}
	if u.LastFailedAt == nil {
		t.Error("expected last_failed_at to be set")
	}
}

func TestFailedAttemptsResetOnSuccess(t *testing.T) {
	s := testStore(t)
	s.CreateUser("admin", "pass", RoleAdmin)

	s.Authenticate("admin", "wrong")
	s.Authenticate("admin", "wrong")

	u, _ := s.GetUser(1)
	if u.FailedAttempts != 2 {
		t.Errorf("failed_attempts = %d, want 2", u.FailedAttempts)
	}

	// Successful login resets counter.
	_, err := s.Authenticate("admin", "pass")
	if err != nil {
		t.Fatal(err)
	}
	u, _ = s.GetUser(1)
	if u.FailedAttempts != 0 {
		t.Errorf("failed_attempts = %d, want 0 after success", u.FailedAttempts)
	}
}

func TestAuthenticateNoUser(t *testing.T) {
	s := testStore(t)
	_, err := s.Authenticate("nonexistent", "pass")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestChangePassword(t *testing.T) {
	s := testStore(t)
	u, _ := s.CreateUser("admin", "old", RoleAdmin)

	if err := s.ChangePassword(u.ID, "old", "new"); err != nil {
		t.Fatal(err)
	}

	_, err := s.Authenticate("admin", "new")
	if err != nil {
		t.Fatal("should authenticate with new password")
	}

	_, err = s.Authenticate("admin", "old")
	if err == nil {
		t.Fatal("should not authenticate with old password")
	}
}

func TestUserCount(t *testing.T) {
	s := testStore(t)

	count, err := s.UserCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}

	s.CreateUser("a", "p", RoleViewer)
	s.CreateUser("b", "p", RoleViewer)

	count, err = s.UserCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestListUsers(t *testing.T) {
	s := testStore(t)
	s.CreateUser("alice", "p", RoleAdmin)
	s.CreateUser("bob", "p", RoleViewer)

	users, err := s.ListUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("len = %d, want 2", len(users))
	}
	if users[0].Username != "alice" || users[0].Role != RoleAdmin {
		t.Errorf("users[0] = %+v", users[0])
	}
	if users[1].Username != "bob" || users[1].Role != RoleViewer {
		t.Errorf("users[1] = %+v", users[1])
	}
}

func TestDeleteUser(t *testing.T) {
	s := testStore(t)
	s.CreateUser("first", "p", RoleAdmin) // id=1
	u, _ := s.CreateUser("victim", "p", RoleViewer)

	if err := s.DeleteUser(u.ID); err != nil {
		t.Fatal(err)
	}

	count, _ := s.UserCount()
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
}

func TestDeleteDefaultAdminBlocked(t *testing.T) {
	s := testStore(t)
	s.EnsureDefaultAdmin() // creates admin with id=1

	err := s.DeleteUser(1)
	if err == nil {
		t.Fatal("expected error when deleting default admin")
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	s := testStore(t)
	err := s.DeleteUser(999)
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestEnsureDefaultAdmin(t *testing.T) {
	s := testStore(t)

	u, err := s.EnsureDefaultAdmin()
	if err != nil {
		t.Fatal(err)
	}
	if u == nil {
		t.Fatal("expected user to be created")
	}
	if u.Username != "admin" || u.Role != RoleAdmin {
		t.Errorf("user = %+v, want admin/admin role", u)
	}

	// Second call should be a no-op.
	u2, err := s.EnsureDefaultAdmin()
	if err != nil {
		t.Fatal(err)
	}
	if u2 != nil {
		t.Error("expected nil on second call (users already exist)")
	}

	if !u.MustChangePassword {
		t.Error("EnsureDefaultAdmin: expected MustChangePassword=true")
	}

	// Verify we can authenticate with admin/admin.
	authed, err := s.Authenticate("admin", "admin")
	if err != nil {
		t.Fatal(err)
	}
	if authed.Role != RoleAdmin {
		t.Errorf("role = %q, want admin", authed.Role)
	}
	if !authed.MustChangePassword {
		t.Error("Authenticate: expected MustChangePassword=true for default admin")
	}

	// Verify GetUser also returns the flag.
	got, err := s.GetUser(authed.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !got.MustChangePassword {
		t.Error("GetUser: expected MustChangePassword=true for default admin")
	}

	// After changing password, flag should be cleared.
	if err := s.ChangePassword(authed.ID, "admin", "newsecret"); err != nil {
		t.Fatal(err)
	}
	got, _ = s.GetUser(authed.ID)
	if got.MustChangePassword {
		t.Error("expected MustChangePassword=false after password change")
	}
}
