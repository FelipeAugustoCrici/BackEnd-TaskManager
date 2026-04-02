-- Cria tabela para armazenar o timer ativo por usuário
CREATE TABLE IF NOT EXISTS active_timers (
  id               UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id          UUID        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  status           VARCHAR(10) NOT NULL CHECK (status IN ('running', 'paused')),
  started_at       TIMESTAMPTZ,
  elapsed_seconds  INT         NOT NULL DEFAULT 0,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_active_timers_user_id ON active_timers (user_id);
