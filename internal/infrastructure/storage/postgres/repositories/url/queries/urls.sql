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

-- name: GetUrlStats :many
-- Получение статистики переходов по короткой ссылке
SELECT
    us.*,
    u.alias,
    u.url
FROM url_stats us
JOIN urls u ON us.url_id = u.id
WHERE u.alias = $1
ORDER BY us.created_at DESC;

-- name: GetUrlStatsByTimeRange :many
-- Статистика с фильтрацией по временному диапазону
SELECT
    us.*,
    u.alias,
    u.url
FROM url_stats us
JOIN urls u ON us.url_id = u.id
WHERE u.alias = $1
AND us.created_at >= $2
AND us.created_at <= $3
ORDER BY us.created_at DESC;

-- name: GetUrlStatsAggregated :many
-- Агрегированная статистика по дням/месяцам и User-Agent
SELECT
    DATE(us.created_at) as date,
    us.user_agent,
    COUNT(*) as visit_count
FROM url_stats us
JOIN urls u ON us.url_id = u.id
WHERE u.alias = $1
AND us.created_at >= $2
AND us.created_at <= $3
GROUP BY DATE(us.created_at), us.user_agent
ORDER BY date DESC, visit_count DESC;

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

-- name: GetURLByID :one
-- Получение URL по ID
SELECT * FROM urls
WHERE id = $1;

-- name: GetTotalVisits :one
-- Общее количество переходов по короткой ссылке
SELECT COUNT(*) as total_visits
FROM url_stats us
JOIN urls u ON us.url_id = u.id
WHERE u.alias = $1;

-- name: GetUniqueVisitors :one
-- Количество уникальных посетителей
SELECT COUNT(DISTINCT user_id) as unique_visitors
FROM url_stats us
JOIN urls u ON us.url_id = u.id
WHERE u.alias = $1;

-- name: GetPopularUserAgents :many
-- Самые популярные User-Agents
SELECT
    user_agent,
    COUNT(*) as count
FROM url_stats us
JOIN urls u ON us.url_id = u.id
WHERE u.alias = $1
GROUP BY user_agent
ORDER BY count DESC
LIMIT 10;