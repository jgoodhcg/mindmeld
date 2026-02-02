---
title: "Question System Enhancements"
status: idea
description: "Expand question types, pools, and AI-assisted creation for trivia."
tags: [area/game, type/feature, tech/ai]
priority: medium
created: 2026-01-19
updated: 2026-02-02
effort: L
depends-on: []
---

# Question System Enhancements

> **NOTE:** AI Question Assistance (MVP) is currently being implemented as part of **[Work Party Prep Phase 2](./work-party-prep.md)**.

## Work Unit Summary

**Problem/Intent:**
Expand beyond "each player writes one multiple-choice question" to support varied question types, reusable question pools, and AI-generated content. Goal: reduce friction for quick games while keeping the personal touch that makes Mindmeld social.

**North Star Check:**
- AI questions should *enhance* conversation, not replace the personal element
- Question pools let groups skip the "everyone write a question" phase when they just want to play
- Variety in question types keeps things fresh and creates different kinds of discussions

---

## Question Types

| Type | Format | Why It's Fun |
|------|--------|--------------|
| **Multiple choice** (current) | 4 options, 1 correct | Classic, low friction |
| **True/False** | 2 options | Fast rounds, hot takes |
| **Free text** | Type answer, fuzzy match | "How would YOU answer this?" |
| **Numeric** | Closest guess wins | Estimation debates |
| **Ranking** | Order items 1-4 | Reveals priorities/values |
| **Fill in the blank** | Complete the phrase | Creative, funny |
| **This or That** | Binary choice, no right answer | Personality reveals |
| **Hot Take** | Agree/disagree scale | Sparks debate |

### Implementation Notes
- Start with True/False (simplest addition)
- Free text needs fuzzy matching (Levenshtein or AI similarity)
- Some types have no "correct" answer — scoring based on matching others (Wavelength-style)

---

## Game Settings

### Round Configuration
| Setting | Options | Default |
|---------|---------|---------|
| **Questions per player** | 1, 2, 3, unlimited | 1 |
| **Total questions** | Fixed number or "all submitted" | All submitted |
| **Question source** | Players only, Pool only, AI only, Mix | Players only |
| **Time limit per question** | None, 15s, 30s, 60s | None |

### Round Styles
| Style | Description |
|-------|-------------|
| **Classic** | Everyone answers, reveal after all submit |
| **Speed** | Timer, fastest correct answer gets bonus |
| **Elimination** | Wrong answer = out, last one standing wins |
| **Team** | Split into teams, team score matters |
| **Head-to-head** | Two players compete per question |

### Game Modes
| Mode | Description |
|------|-------------|
| **Party** (current) | Everyone submits, everyone plays |
| **Quiz Bowl** | One person/team runs questions, others compete |
| **AI Challenge** | AI generates all questions based on theme |
| **Mixed** | Some player questions, some AI, some from pool |

---

## Question Sources

### 1. Player-Submitted (Current)
- Each player writes during "submitting" phase
- Personal, creates "who wrote this?" moments
- Requires everyone to participate upfront

### 2. Question Pool
**Concept:** Previously submitted questions can be saved and reused.

| Feature | Description |
|---------|-------------|
| **Opt-in saving** | "Add this question to the community pool?" |
| **Categories/tags** | Pop culture, science, personal, etc. |
| **Difficulty rating** | Derived from answer accuracy over time |
| **Attribution** | Optional: "Question by @username" |
| **Flagging** | Report inappropriate questions |
| **Private pools** | Groups can build their own pool over time |

**Use Cases:**
- Quick game: "Just give us 10 random questions"
- Themed party: "Only 90s pop culture"
- Replay value: Questions you haven't seen before

### 3. AI-Assisted Question Creation

**Concept:** AI helps players write better questions faster, not replace the personal touch.

#### Assistance Modes

**A. Trivia Category Questions**
- User picks a category (Animals, History, Pop Culture, etc.)
- AI generates a complete question with correct + wrong answers
- Example: "Which of these is NOT a reptile?" → Correct: Salamander, Wrong: Gecko, Iguana, Turtle
- User can edit/tweak before submitting

**B. Personal Questions About Players**
- User picks a player or writes "[Name]'s..."
- AI generates question template: "What is [Name]'s favorite color?"
- **Correct answer: Left blank** (only that person knows)
- **Wrong answers: Optional AI assist** (user can request plausible wrong answers or write their own)
- Harder for AI to know personal facts, so default to manual entry

#### What AI Can Help With

| Task | Trivia Mode | Personal Mode |
|------|-------------|---------------|
| Generate question | ✅ Full auto | ✅ Template only |
| Fill correct answer | ✅ Full auto | ❌ User must fill |
| Fill wrong answers | ✅ Full auto | ⚠️ Optional (user can request) |
| Suggest categories | ✅ | ❌ |

#### UI Flow

```
[Write Question]  [AI Assist ▼]
                    ├── Trivia Category → Pick category → Generate
                    └── About a Player → Pick player → Generate template

[Question text field - editable]

Correct Answer: [field]  [🤖 Suggest]
Wrong Answer 1: [field]  [🤖 Suggest]
Wrong Answer 2: [field]  [🤖 Suggest]
Wrong Answer 3: [field]  [🤖 Suggest]

[Submit Question]
```

---

### 4. Full AI Game Mode (Future)

For when you want to skip question writing entirely:

| Setting | Options |
|---------|---------|
| **Theme** | General, Pop Culture, Science, History, Sports, Custom... |
| **Difficulty** | Easy, Medium, Hard, Expert |
| **Demographics** | Millennials, Gen Z, Boomers, Kids, Mixed |
| **Tone** | Serious, Silly, Weird, Educational |
| **Question count** | 5, 10, 15, 20 |
| **Custom prompt** | "Questions about 80s action movies" |

**Hybrid Modes:**
- "3 player questions + 2 AI questions per round"
- "AI writes questions, players write wrong answers" (Fibbage-style)

---

## AI Cost & Access Model

**The Problem:** AI API calls cost money. Who pays?

### Phase 1: Platform-Funded (MVP)
- Use platform API keys
- Use cheap/fast models (GPT-4o-mini, Claude Haiku, Gemini Flash)
- Rate limit per user/session (e.g., 10 AI assists per game)
- Monitor costs, adjust limits as needed
- Risk: Could get expensive if it goes viral

### Phase 2: Accounts Required
- AI features require login
- Better rate limiting per account
- Can track abuse, ban bad actors
- Still platform-funded but controlled

### Phase 3: Bring Your Own Key (BYOK)
- Users can add their own API keys in settings
- Unlocks unlimited AI usage
- Can choose their preferred model
- Platform-funded tier still available with limits

### Model Tiers

| Tier | Models | Cost | Who Pays |
|------|--------|------|----------|
| **Free (limited)** | Haiku, Gemini Flash, GPT-4o-mini | $ | Platform |
| **Free (logged in)** | Same, higher limits | $ | Platform |
| **BYOK** | Any model user wants | $$$ | User |

### Rate Limiting Ideas
- X AI assists per game (resets each game)
- X AI assists per day (resets daily)
- Cooldown between requests (prevent spam)
- Higher limits for logged-in users
- Unlimited for BYOK users

---

## Data Model Additions

```sql
-- Question types
ALTER TABLE trivia_questions ADD COLUMN question_type VARCHAR DEFAULT 'multiple_choice';
-- Values: multiple_choice, true_false, free_text, numeric, ranking, fill_blank, this_or_that, hot_take

-- Question pool
CREATE TABLE question_pool (
  id UUID PRIMARY KEY,
  question_text TEXT,
  question_type VARCHAR,
  correct_answer VARCHAR,
  wrong_answers JSONB,  -- Flexible for different types
  category VARCHAR,
  difficulty_rating FLOAT,  -- Computed from usage
  times_used INT DEFAULT 0,
  times_correct INT DEFAULT 0,
  created_by UUID REFERENCES players(id),
  is_ai_generated BOOLEAN DEFAULT FALSE,
  ai_prompt TEXT,  -- What prompt generated this
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Game configuration
CREATE TABLE game_configs (
  id UUID PRIMARY KEY,
  lobby_id UUID REFERENCES lobbies(id),
  questions_per_player INT DEFAULT 1,
  question_source VARCHAR DEFAULT 'players', -- players, pool, ai, mixed
  time_limit_seconds INT,
  round_style VARCHAR DEFAULT 'classic',
  ai_theme VARCHAR,
  ai_difficulty VARCHAR,
  ai_demographics VARCHAR,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Prioritization (WSJF)

| # | Feature | Value | Size | Priority |
|---|---------|-------|------|----------|
| 1 | True/False question type | Medium | Small | High |
| 2 | AI: Generate trivia question (category-based) | High | Medium | High |
| 3 | AI: Generate personal question template | Medium | Small | High |
| 4 | AI: Suggest wrong answers | Medium | Small | High |
| 5 | Questions per player setting | Medium | Small | Medium |
| 6 | Time limit option | Medium | Small | Medium |
| 7 | Accounts system (for AI rate limiting) | High | Large | Medium |
| 8 | Question pool (save/reuse) | Medium | Large | Medium |
| 9 | BYOK (bring your own API key) | Medium | Medium | Low |
| 10 | Full AI game mode | Medium | Large | Low |
| 11 | Free text with fuzzy match | Medium | Medium | Low |
| 12 | Additional question types | Low | Medium | Low |

---

## Open Questions

- Should AI-assisted questions be visually distinct? ("AI-assisted" badge)
- How to handle inappropriate AI outputs? (Moderation layer, content filtering)
- Private vs public question pools — privacy implications?
- How does AI question generation fit the "promote conversation" north star?
- What's the right rate limit for platform-funded AI? (cost vs friction)
- Should personal question templates know player names from the lobby?
- When to require accounts? (AI features only, or eventually everything?)
- How to prevent abuse of platform-funded AI tier?
