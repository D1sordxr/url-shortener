-- name: GetURLByID :one
-- Получение URL по ID
SELECT * FROM urls
WHERE id = $1;


-- name: CreateURL :one
-- Запрос создает url
-- Возвращает созданную запись целиком
INSERT INTO urls (
    alias,
    url
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetByAlias :one
-- Запрос получает url по alias
SELECT * FROM urls
WHERE alias = $1;

-- name: GetCompleteAnalytics :one
WITH stats AS (
    SELECT
        us.url_id,
        COUNT(*) as total_visits,
        COUNT(DISTINCT us.user_id) as unique_visitors,
        MIN(us.created_at) as first_visit,
        MAX(us.created_at) as last_visit
    FROM url_stats us
    GROUP BY us.url_id
),
     recent_visits AS (
         SELECT
             url_id,
             JSON_AGG(
                     JSON_BUILD_OBJECT(
                             'date', DATE(created_at),
                             'user_agent', user_agent,
                             'ip_address', ip_address,
                             'referer', referer,
                             'created_at', created_at
                     ) ORDER BY created_at DESC
             ) as raw_stats
         FROM url_stats
         WHERE url_id IN (SELECT id FROM urls WHERE alias = $1)
         GROUP BY url_id
     )
SELECT
    u.alias,
    u.url as original_url,
    COALESCE(s.total_visits, 0) as total_visits,
    COALESCE(s.unique_visitors, 0) as unique_visitors,
    s.first_visit,
    s.last_visit,
    COALESCE(rv.raw_stats, '[]'::json) as raw_stats
FROM urls u
         LEFT JOIN stats s ON u.id = s.url_id
         LEFT JOIN recent_visits rv ON u.id = rv.url_id
WHERE u.alias = $1;

-- name: CreateURLStat :one
-- Создание записи статистики перехода
INSERT INTO url_stats (
    url_id,
    user_id,
    user_agent,
    ip_address,
    referer
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

