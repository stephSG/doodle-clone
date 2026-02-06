package database

func createUsersTable() string {
	return `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255),
		name VARCHAR(255) NOT NULL,
		avatar VARCHAR(500),
		provider VARCHAR(50) NOT NULL DEFAULT 'email',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
}

func createPollsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS polls (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(200) NOT NULL,
		description TEXT,
		location VARCHAR(500),
		creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		expires_at TIMESTAMP WITH TIME ZONE,
		allow_multiple BOOLEAN DEFAULT false,
		allow_maybe BOOLEAN DEFAULT true,
		anonymous BOOLEAN DEFAULT false,
		limit_votes BOOLEAN DEFAULT false,
		max_votes_per_user INTEGER DEFAULT 1,
		final_date UUID REFERENCES date_options(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_polls_creator ON polls(creator_id);
	CREATE INDEX IF NOT EXISTS idx_polls_expires ON polls(expires_at);
	CREATE INDEX IF NOT EXISTS idx_polls_final_date ON polls(final_date);
	`
}

func createDateOptionsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS date_options (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
		start_time TIMESTAMP WITH TIME ZONE NOT NULL,
		end_time TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_date_options_poll ON date_options(poll_id);
	`
}

func createVotesTable() string {
	return `
	CREATE TABLE IF NOT EXISTS votes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
		date_option_id UUID NOT NULL REFERENCES date_options(id) ON DELETE CASCADE,
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		user_name VARCHAR(255) NOT NULL,
		response VARCHAR(10) NOT NULL CHECK (response IN ('yes', 'no', 'maybe')),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(poll_id, date_option_id, user_id)
	);

	CREATE INDEX IF NOT EXISTS idx_votes_poll ON votes(poll_id);
	CREATE INDEX IF NOT EXISTS idx_votes_date_option ON votes(date_option_id);
	CREATE INDEX IF NOT EXISTS idx_votes_user ON votes(user_id);
	CREATE INDEX IF NOT EXISTS idx_votes_poll_user ON votes(poll_id, user_id);
	`
}

func createCommentsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS comments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_comments_poll ON comments(poll_id);
	CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id);
	`
}

func createRefreshTokensTable() string {
	return `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(500) UNIQUE NOT NULL,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at);
	`
}

func createNotificationSettingsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS notification_settings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		key VARCHAR(100) UNIQUE NOT NULL,
		value TEXT NOT NULL,
		description TEXT,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
}

func createNotificationsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		type VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
		sent_at TIMESTAMP WITH TIME ZONE,
		error_message TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_notifications_poll ON notifications(poll_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
	CREATE INDEX IF NOT EXISTS idx_notifications_scheduled ON notifications(scheduled_at);
	`
}
