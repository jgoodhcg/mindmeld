-- +goose Up

-- Retire original starter content (too small / too mild for game-day sessions).
UPDATE coordinates_prompt_axis_sets
SET is_active = FALSE
WHERE id IN (
    '33333333-3333-3333-3333-333333333301',
    '33333333-3333-3333-3333-333333333302',
    '33333333-3333-3333-3333-333333333303'
);

UPDATE coordinates_prompts
SET is_active = FALSE
WHERE created_by_label = 'cluster-seed-v1';

UPDATE coordinates_axis_sets
SET is_active = FALSE
WHERE created_by_label = 'cluster-seed-v1';

INSERT INTO coordinates_axis_sets (
    id,
    x_min_label,
    x_max_label,
    y_min_label,
    y_max_label,
    created_by_kind,
    created_by_label,
    authoring_mode,
    generator_provider,
    generator_model,
    generator_prompt_version,
    generator_run_id,
    provenance,
    min_rating,
    is_active
) VALUES
(
    '44444444-4444-4444-4444-444444444401',
    'Consensus-first',
    'Contrarian',
    'Low stakes',
    'High stakes',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","note":"cluster v2 axis"}'::jsonb,
    20,
    TRUE
),
(
    '44444444-4444-4444-4444-444444444402',
    'Personal preference',
    'Team impact',
    'Fast decision',
    'Deliberate decision',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","note":"cluster v2 axis"}'::jsonb,
    20,
    TRUE
),
(
    '44444444-4444-4444-4444-444444444403',
    'Async-friendly',
    'Live collaboration',
    'Structured',
    'Flexible',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","note":"cluster v2 axis"}'::jsonb,
    20,
    TRUE
),
(
    '44444444-4444-4444-4444-444444444404',
    'Budget-friendly',
    'Premium',
    'Calm energy',
    'High energy',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","note":"cluster v2 axis"}'::jsonb,
    10,
    TRUE
),
(
    '44444444-4444-4444-4444-444444444405',
    'Classic approach',
    'Experimental approach',
    'Independent work',
    'Pair or group work',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","note":"cluster v2 axis"}'::jsonb,
    10,
    TRUE
),
(
    '44444444-4444-4444-4444-444444444406',
    'Practical',
    'Aspirational',
    'Immediate payoff',
    'Long-term payoff',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","note":"cluster v2 axis"}'::jsonb,
    20,
    TRUE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO coordinates_prompts (
    id,
    prompt_text,
    created_by_kind,
    created_by_label,
    authoring_mode,
    generator_provider,
    generator_model,
    generator_prompt_version,
    generator_run_id,
    provenance,
    min_rating,
    is_active
) VALUES
(
    '55555555-5555-5555-5555-555555555501',
    'The best way to kick off a Monday team sync',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"meetings"}'::jsonb,
    20,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555502',
    'An ideal format for sharing project updates',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"communication"}'::jsonb,
    20,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555503',
    'The most useful retro activity after a rough sprint',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"retrospectives"}'::jsonb,
    20,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555504',
    'The best way to celebrate shipping a big feature',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"celebration"}'::jsonb,
    10,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555505',
    'The strongest signal a meeting should have been async',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"focus time"}'::jsonb,
    20,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555506',
    'A great onboarding moment for a new teammate',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"onboarding"}'::jsonb,
    20,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555507',
    'The most effective way to resolve cross-team disagreements',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"alignment"}'::jsonb,
    20,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555508',
    'The best 15-minute reset during a stressful day',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"energy management"}'::jsonb,
    10,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555509',
    'A team ritual worth protecting as the company grows',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"team culture"}'::jsonb,
    20,
    TRUE
),
(
    '55555555-5555-5555-5555-555555555510',
    'The most motivating recognition for great work',
    'developer',
    'cluster-seed-v2',
    'manual',
    NULL,
    NULL,
    NULL,
    NULL,
    '{"seeded_by":"developer","topic":"recognition"}'::jsonb,
    10,
    TRUE
)
ON CONFLICT (id) DO NOTHING;

-- Cross-join all v2 prompts and axis sets: 10 x 6 = 60 combinations.
INSERT INTO coordinates_prompt_axis_sets (
    id,
    prompt_id,
    axis_set_id,
    is_active
)
SELECT
    gen_random_uuid(),
    p.id,
    a.id,
    TRUE
FROM coordinates_prompts p
JOIN coordinates_axis_sets a
    ON a.created_by_label = 'cluster-seed-v2'
WHERE p.created_by_label = 'cluster-seed-v2'
ON CONFLICT (prompt_id, axis_set_id) DO NOTHING;

-- +goose Down

DELETE FROM coordinates_prompt_axis_sets cpas
USING coordinates_prompts cp, coordinates_axis_sets cas
WHERE cpas.prompt_id = cp.id
  AND cpas.axis_set_id = cas.id
  AND cp.created_by_label = 'cluster-seed-v2'
  AND cas.created_by_label = 'cluster-seed-v2';

DELETE FROM coordinates_prompts
WHERE created_by_label = 'cluster-seed-v2';

DELETE FROM coordinates_axis_sets
WHERE created_by_label = 'cluster-seed-v2';

UPDATE coordinates_prompt_axis_sets
SET is_active = TRUE
WHERE id IN (
    '33333333-3333-3333-3333-333333333301',
    '33333333-3333-3333-3333-333333333302',
    '33333333-3333-3333-3333-333333333303'
);

UPDATE coordinates_prompts
SET is_active = TRUE
WHERE created_by_label = 'cluster-seed-v1';

UPDATE coordinates_axis_sets
SET is_active = TRUE
WHERE created_by_label = 'cluster-seed-v1';
