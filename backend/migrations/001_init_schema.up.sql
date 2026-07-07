-- ============================================
-- 001_init_schema.up.sql
-- 大学生就业需求分析系统 - 初始化建表
-- PostgreSQL 16
-- ============================================

BEGIN;

-- 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- 1. 系统基础表
-- ============================================

CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(64)  NOT NULL UNIQUE,
    password_hash   VARCHAR(256) NOT NULL,
    nickname        VARCHAR(64),
    email           BYTEA,
    phone           BYTEA,
    avatar          VARCHAR(512),
    status          SMALLINT     NOT NULL DEFAULT 1,
    login_attempts  SMALLINT     NOT NULL DEFAULT 0,
    locked_until    TIMESTAMPTZ,
    last_login_at   TIMESTAMPTZ,
    last_login_ip   VARCHAR(45),
    password_changed_at TIMESTAMPTZ,
    must_change_password BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted ON users(deleted_at);

CREATE TABLE roles (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(64)  NOT NULL UNIQUE,
    code        VARCHAR(64)  NOT NULL UNIQUE,
    description VARCHAR(256),
    status      SMALLINT     NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE user_roles (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE menus (
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGINT       REFERENCES menus(id),
    name        VARCHAR(64)  NOT NULL,
    path        VARCHAR(256),
    component   VARCHAR(256),
    icon        VARCHAR(64),
    sort_order  INT          NOT NULL DEFAULT 0,
    permission  VARCHAR(128),
    menu_type   SMALLINT     NOT NULL DEFAULT 2,
    status      SMALLINT     NOT NULL DEFAULT 1,
    visible     BOOLEAN      NOT NULL DEFAULT TRUE,
    keep_alive  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE role_menus (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    menu_id BIGINT NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, menu_id)
);

CREATE TABLE dict_data (
    id          BIGSERIAL PRIMARY KEY,
    dict_type   VARCHAR(64)  NOT NULL,
    dict_label  VARCHAR(128) NOT NULL,
    dict_value  VARCHAR(128) NOT NULL,
    sort_order  INT          NOT NULL DEFAULT 0,
    status      SMALLINT     NOT NULL DEFAULT 1,
    remark      VARCHAR(256),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(dict_type, dict_value)
);
CREATE INDEX idx_dict_type ON dict_data(dict_type);

-- ============================================
-- 2. 学院与专业
-- ============================================

CREATE TABLE colleges (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(128) NOT NULL,
    code       VARCHAR(32)  UNIQUE,
    dean_name  VARCHAR(64),
    contact    VARCHAR(128),
    sort_order INT          NOT NULL DEFAULT 0,
    status     SMALLINT     NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE majors (
    id              BIGSERIAL PRIMARY KEY,
    college_id      BIGINT       NOT NULL REFERENCES colleges(id),
    name            VARCHAR(128) NOT NULL,
    code            VARCHAR(32)  UNIQUE,
    category        VARCHAR(32),
    education_level VARCHAR(16),
    duration_years  SMALLINT DEFAULT 4,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_majors_college ON majors(college_id);

-- ============================================
-- 3. 学生与技能
-- ============================================

CREATE TABLE students (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    student_no          VARCHAR(32)  NOT NULL UNIQUE,
    name                VARCHAR(64)  NOT NULL,
    gender              SMALLINT,
    birth_date          DATE,
    college_id          BIGINT       REFERENCES colleges(id),
    major_id            BIGINT       REFERENCES majors(id),
    grade               VARCHAR(4),
    education_level     VARCHAR(16),
    graduation_year     INT,
    political_status    VARCHAR(16),
    hometown_city       VARCHAR(64),
    phone               BYTEA,
    email               BYTEA,
    wechat              VARCHAR(64),
    qq                  VARCHAR(20),
    id_card_hash        VARCHAR(64),
    resume_url          VARCHAR(512),
    expected_city       VARCHAR(64),
    expected_salary_min INT,
    expected_salary_max INT,
    expected_industry   VARCHAR(64),
    expected_positions  TEXT,
    employ_status       VARCHAR(16) DEFAULT 'unemployed',
    employ_company      VARCHAR(128),
    employ_position     VARCHAR(128),
    employ_salary       INT,
    employ_city         VARCHAR(64),
    employ_date         DATE,
    data_source         VARCHAR(16) DEFAULT 'manual',
    verified            BOOLEAN DEFAULT FALSE,
    remark              TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);
CREATE INDEX idx_students_college     ON students(college_id);
CREATE INDEX idx_students_major       ON students(major_id);
CREATE INDEX idx_students_employ      ON students(employ_status);
CREATE INDEX idx_students_grad_year   ON students(graduation_year);
CREATE INDEX idx_students_deleted     ON students(deleted_at);

CREATE TABLE skills (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(64)  NOT NULL UNIQUE,
    category   VARCHAR(32),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE student_skills (
    student_id  BIGINT  NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    skill_id    BIGINT  NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    proficiency SMALLINT DEFAULT 1 CHECK (proficiency BETWEEN 1 AND 5),
    PRIMARY KEY (student_id, skill_id)
);

-- ============================================
-- 4. 行业与企业
-- ============================================

CREATE TABLE industries (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(64)  NOT NULL UNIQUE,
    code       VARCHAR(32)  UNIQUE,
    parent_id  BIGINT       REFERENCES industries(id),
    sort_order INT DEFAULT 0,
    status     SMALLINT DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE enterprises (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    name            VARCHAR(128) NOT NULL,
    short_name      VARCHAR(64),
    industry_id     BIGINT       REFERENCES industries(id),
    scale           VARCHAR(16),
    nature          VARCHAR(16),
    city            VARCHAR(64),
    address         VARCHAR(256),
    website         VARCHAR(256),
    logo_url        VARCHAR(512),
    description     TEXT,
    business_scope  TEXT,
    contact_name    VARCHAR(64),
    contact_phone   BYTEA,
    contact_email   BYTEA,
    contact_position VARCHAR(64),
    status          SMALLINT DEFAULT 1,
    verified        BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_enterprises_industry ON enterprises(industry_id);

-- ============================================
-- 5. 职位
-- ============================================

CREATE TABLE positions (
    id                     BIGSERIAL PRIMARY KEY,
    enterprise_id          BIGINT       NOT NULL REFERENCES enterprises(id),
    title                  VARCHAR(128) NOT NULL,
    industry_id            BIGINT       REFERENCES industries(id),
    city                   VARCHAR(64),
    district               VARCHAR(64),
    address                VARCHAR(256),
    education_requirement  VARCHAR(16),
    experience_requirement INT,
    salary_min             INT,
    salary_max             INT,
    salary_type            VARCHAR(16) DEFAULT 'monthly',
    headcount              INT DEFAULT 1,
    job_type               VARCHAR(16) DEFAULT 'fulltime',
    description            TEXT,
    requirement            TEXT,
    welfare                TEXT,
    contact_info           VARCHAR(256),
    status                 SMALLINT DEFAULT 1,
    publish_at             TIMESTAMPTZ,
    expire_at              TIMESTAMPTZ,
    view_count             INT DEFAULT 0,
    apply_count            INT DEFAULT 0,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ
);
CREATE INDEX idx_positions_enterprise ON positions(enterprise_id);
CREATE INDEX idx_positions_industry   ON positions(industry_id);
CREATE INDEX idx_positions_status     ON positions(status);

CREATE TABLE position_skills (
    position_id BIGINT  NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
    skill_id    BIGINT  NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    importance  SMALLINT DEFAULT 1 CHECK (importance BETWEEN 1 AND 5),
    PRIMARY KEY (position_id, skill_id)
);

-- ============================================
-- 6. 求职投递
-- ============================================

CREATE TABLE applications (
    id             BIGSERIAL PRIMARY KEY,
    student_id     BIGINT       NOT NULL REFERENCES students(id),
    position_id    BIGINT       NOT NULL REFERENCES positions(id),
    enterprise_id  BIGINT       NOT NULL REFERENCES enterprises(id),
    status         VARCHAR(16)  DEFAULT 'pending',
    resume_url     VARCHAR(512),
    cover_letter   TEXT,
    interview_at   TIMESTAMPTZ,
    interview_note TEXT,
    offer_salary   INT,
    offer_at       TIMESTAMPTZ,
    reject_reason  VARCHAR(256),
    remark         TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_position UNIQUE(student_id, position_id)
);
CREATE INDEX idx_applications_student   ON applications(student_id);
CREATE INDEX idx_applications_position  ON applications(position_id);
CREATE INDEX idx_applications_status    ON applications(status);

CREATE TABLE application_logs (
    id             BIGSERIAL PRIMARY KEY,
    application_id BIGINT       NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    from_status    VARCHAR(16),
    to_status      VARCHAR(16)  NOT NULL,
    operator_id    BIGINT       REFERENCES users(id),
    operator_name  VARCHAR(64),
    note           TEXT,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================
-- 7. 就业记录
-- ============================================

CREATE TABLE employment_records (
    id               BIGSERIAL PRIMARY KEY,
    student_id       BIGINT       UNIQUE REFERENCES students(id),
    company_name     VARCHAR(128) NOT NULL,
    position_name    VARCHAR(128),
    industry_id      BIGINT       REFERENCES industries(id),
    city             VARCHAR(64),
    monthly_salary   INT,
    annual_salary    INT,
    contract_type    VARCHAR(32),
    employment_date  DATE,
    probation_months INT,
    social_security  BOOLEAN,
    data_source      VARCHAR(16) DEFAULT 'manual',
    verified         BOOLEAN DEFAULT FALSE,
    verified_by      BIGINT REFERENCES users(id),
    verified_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- 8. 就业调研
-- ============================================

CREATE TABLE employment_surveys (
    id               BIGSERIAL PRIMARY KEY,
    student_id       BIGINT       REFERENCES students(id),
    survey_year      INT          NOT NULL,
    survey_type      VARCHAR(16) DEFAULT 'graduate',
    employ_status    VARCHAR(16),
    company_name     VARCHAR(128),
    industry_id      BIGINT       REFERENCES industries(id),
    city             VARCHAR(64),
    monthly_salary   INT,
    position_name    VARCHAR(128),
    satisfaction     SMALLINT     CHECK (satisfaction BETWEEN 1 AND 5),
    major_related    BOOLEAN,
    skill_match      SMALLINT     CHECK (skill_match BETWEEN 1 AND 5),
    job_source       VARCHAR(64),
    feedback         TEXT,
    submitted_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(student_id, survey_year, survey_type)
);

-- ============================================
-- 9. 分析报告 & 市场数据
-- ============================================

CREATE TABLE analysis_reports (
    id             BIGSERIAL PRIMARY KEY,
    title          VARCHAR(256) NOT NULL,
    report_type    VARCHAR(32)  NOT NULL,
    report_scope   VARCHAR(16) DEFAULT 'university',
    scope_id       BIGINT,
    report_year    INT,
    report_period  VARCHAR(16),
    content        JSONB,
    summary        TEXT,
    file_url       VARCHAR(512),
    file_format    VARCHAR(16),
    status         VARCHAR(16) DEFAULT 'draft',
    created_by     BIGINT REFERENCES users(id),
    published_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE market_data_snapshots (
    id              BIGSERIAL PRIMARY KEY,
    source          VARCHAR(128),
    industry_id     BIGINT REFERENCES industries(id),
    city            VARCHAR(64),
    position_name   VARCHAR(128),
    avg_salary      INT,
    salary_p25      INT,
    salary_p50      INT,
    salary_p75      INT,
    demand_count    INT,
    supply_count    INT,
    growth_rate     DECIMAL(5,2),
    snapshot_date   DATE NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, industry_id, city, position_name, snapshot_date)
);

-- ============================================
-- 10. 操作日志（不可物理删除）
-- ============================================

CREATE TABLE operation_logs (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT,
    username     VARCHAR(64)  NOT NULL,
    role         VARCHAR(64),
    module       VARCHAR(64)  NOT NULL,
    action       VARCHAR(64)  NOT NULL,
    target_type  VARCHAR(64),
    target_id    BIGINT,
    description  VARCHAR(512),
    method       VARCHAR(10),
    url          VARCHAR(256),
    ip           VARCHAR(45)  NOT NULL,
    user_agent   TEXT,
    request_body TEXT,
    status_code  INT,
    cost_ms      INT,
    is_deleted   BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_operation_logs_created ON operation_logs(created_at DESC);
CREATE INDEX idx_operation_logs_user    ON operation_logs(user_id);

-- 禁止物理删除审计日志
CREATE OR REPLACE FUNCTION prevent_audit_delete()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION '审计日志不可删除！';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_audit_delete
    BEFORE DELETE ON operation_logs
    FOR EACH ROW EXECUTE FUNCTION prevent_audit_delete();

-- 自动更新 updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE t text;
BEGIN
    FOR t IN SELECT table_name FROM information_schema.columns
             WHERE column_name = 'updated_at' AND table_schema = 'public'
    LOOP
        EXECUTE format('CREATE TRIGGER trg_%I_updated_at BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()', t, t);
    END LOOP;
END;
$$;

COMMIT;
