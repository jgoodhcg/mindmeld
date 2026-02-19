package questions

import (
	"sort"

	"github.com/jgoodhcg/mindmeld/internal/contentrating"
)

// Pack represents a curated collection of templates.
type Pack struct {
	ID          string
	Name        string
	Description string
	MinRating   int16
}

// Template represents a ready-to-play trivia question template.
type Template struct {
	ID            string
	PackID        string
	Category      string
	QuestionText  string
	CorrectAnswer string
	WrongAnswer1  string
	WrongAnswer2  string
	WrongAnswer3  string
	MinRating     int16
}

// PackSection is the render-ready shape used by the templates modal.
type PackSection struct {
	Pack                Pack
	Categories          []string
	TemplatesByCategory map[string][]Template
	TemplateCount       int
}

const (
	PackWorkEssentials = "work-essentials"
	PackProductTech    = "product-tech-basics"
	PackOfficePop      = "office-pop-culture"
	PackQuickBrain     = "quick-brain-boost"
	PackWorldSnapshot  = "world-snapshot"
)

// Categories for organizing templates inside each pack.
const (
	CategoryMeetingsProcess = "Meetings & Process"
	CategoryProductDelivery = "Product Delivery"
	CategoryWebFundamentals = "Web Fundamentals"
	CategoryDevWorkflow     = "Dev Workflow"
	CategoryTVFilm          = "TV & Film"
	CategoryMusicCulture    = "Music & Culture"
	CategoryScienceMath     = "Science & Math"
	CategoryEverydayFacts   = "Everyday Facts"
	CategoryHistory         = "History"
	CategoryGeography       = "Geography"
)

var categoryOrder = []string{
	CategoryMeetingsProcess,
	CategoryProductDelivery,
	CategoryWebFundamentals,
	CategoryDevWorkflow,
	CategoryTVFilm,
	CategoryMusicCulture,
	CategoryScienceMath,
	CategoryEverydayFacts,
	CategoryHistory,
	CategoryGeography,
}

// AllPacks contains all curated template packs, ordered for display.
var AllPacks = []Pack{
	{
		ID:          PackWorkEssentials,
		Name:        "Work Essentials",
		Description: "Fast-start work-safe questions about meetings, delivery, and team process.",
		MinRating:   contentrating.Work,
	},
	{
		ID:          PackProductTech,
		Name:        "Product & Tech Basics",
		Description: "Practical software and web fundamentals for mixed technical teams.",
		MinRating:   contentrating.Work,
	},
	{
		ID:          PackOfficePop,
		Name:        "Office Pop Culture",
		Description: "Work-safe pop culture prompts with broad recognition.",
		MinRating:   contentrating.Work,
	},
	{
		ID:          PackQuickBrain,
		Name:        "Quick Brain Boost",
		Description: "Simple general-knowledge questions for all audiences.",
		MinRating:   contentrating.Kids,
	},
	{
		ID:          PackWorldSnapshot,
		Name:        "World Snapshot",
		Description: "Classic history and geography questions that play well in groups.",
		MinRating:   contentrating.Kids,
	},
}

// AllTemplates contains all available templates.
var AllTemplates = []Template{
	// Work Essentials
	{
		ID:            "work-001",
		PackID:        PackWorkEssentials,
		Category:      CategoryMeetingsProcess,
		QuestionText:  `In a RACI matrix, what does the "A" stand for?`,
		CorrectAnswer: "Accountable",
		WrongAnswer1:  "Available",
		WrongAnswer2:  "Approved",
		WrongAnswer3:  "Assigned",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "work-002",
		PackID:        PackWorkEssentials,
		Category:      CategoryProductDelivery,
		QuestionText:  "What does OKR stand for?",
		CorrectAnswer: "Objectives and Key Results",
		WrongAnswer1:  "Operations and Knowledge Review",
		WrongAnswer2:  "Objectives and KPI Reporting",
		WrongAnswer3:  "Outcomes, Knowledge, and Roadmaps",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "work-003",
		PackID:        PackWorkEssentials,
		Category:      CategoryProductDelivery,
		QuestionText:  "In Scrum, who usually prioritizes the product backlog?",
		CorrectAnswer: "Product Owner",
		WrongAnswer1:  "Engineering Manager",
		WrongAnswer2:  "Scrum Master",
		WrongAnswer3:  "QA Lead",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "work-004",
		PackID:        PackWorkEssentials,
		Category:      CategoryMeetingsProcess,
		QuestionText:  "What is the main goal of a sprint retrospective?",
		CorrectAnswer: "Improve how the team works",
		WrongAnswer1:  "Assign performance ratings",
		WrongAnswer2:  "Rewrite the product roadmap",
		WrongAnswer3:  "Choose next sprint's holiday schedule",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "work-005",
		PackID:        PackWorkEssentials,
		Category:      CategoryProductDelivery,
		QuestionText:  "What does MVP stand for in product development?",
		CorrectAnswer: "Minimum Viable Product",
		WrongAnswer1:  "Most Valuable Proposal",
		WrongAnswer2:  "Managed Validation Process",
		WrongAnswer3:  "Minimum Visual Prototype",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "work-006",
		PackID:        PackWorkEssentials,
		Category:      CategoryMeetingsProcess,
		QuestionText:  "In project updates, ETA usually means:",
		CorrectAnswer: "Estimated Time of Arrival",
		WrongAnswer1:  "Estimated Task Assignment",
		WrongAnswer2:  "Expected Team Action",
		WrongAnswer3:  "Effective Turnaround Agreement",
		MinRating:     contentrating.Work,
	},

	// Product & Tech Basics
	{
		ID:            "tech-001",
		PackID:        PackProductTech,
		Category:      CategoryWebFundamentals,
		QuestionText:  `Which HTTP status code means "Not Found"?`,
		CorrectAnswer: "404",
		WrongAnswer1:  "401",
		WrongAnswer2:  "302",
		WrongAnswer3:  "500",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "tech-002",
		PackID:        PackProductTech,
		Category:      CategoryDevWorkflow,
		QuestionText:  "Which Git command creates a local copy of a repository?",
		CorrectAnswer: "git clone",
		WrongAnswer1:  "git fork",
		WrongAnswer2:  "git init",
		WrongAnswer3:  "git copy",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "tech-003",
		PackID:        PackProductTech,
		Category:      CategoryWebFundamentals,
		QuestionText:  "In SQL, which keyword filters rows before grouping?",
		CorrectAnswer: "WHERE",
		WrongAnswer1:  "HAVING",
		WrongAnswer2:  "ORDER BY",
		WrongAnswer3:  "LIMIT",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "tech-004",
		PackID:        PackProductTech,
		Category:      CategoryWebFundamentals,
		QuestionText:  "What does API stand for?",
		CorrectAnswer: "Application Programming Interface",
		WrongAnswer1:  "Automated Program Integration",
		WrongAnswer2:  "Application Process Input",
		WrongAnswer3:  "Advanced Protocol Interface",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "tech-005",
		PackID:        PackProductTech,
		Category:      CategoryDevWorkflow,
		QuestionText:  `In CI/CD, what does "CI" stand for?`,
		CorrectAnswer: "Continuous Integration",
		WrongAnswer1:  "Code Inspection",
		WrongAnswer2:  "Change Implementation",
		WrongAnswer3:  "Continuous Improvement",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "tech-006",
		PackID:        PackProductTech,
		Category:      CategoryWebFundamentals,
		QuestionText:  "Which CSS property controls spacing inside an element's border?",
		CorrectAnswer: "padding",
		WrongAnswer1:  "margin",
		WrongAnswer2:  "gap",
		WrongAnswer3:  "outline",
		MinRating:     contentrating.Work,
	},

	// Office Pop Culture
	{
		ID:            "pop-001",
		PackID:        PackOfficePop,
		Category:      CategoryTVFilm,
		QuestionText:  "Which TV comedy is set at Dunder Mifflin?",
		CorrectAnswer: "The Office",
		WrongAnswer1:  "Parks and Recreation",
		WrongAnswer2:  "Brooklyn Nine-Nine",
		WrongAnswer3:  "Community",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "pop-002",
		PackID:        PackOfficePop,
		Category:      CategoryTVFilm,
		QuestionText:  "In Friends, what is the name of the coffee shop hangout?",
		CorrectAnswer: "Central Perk",
		WrongAnswer1:  "Coffee Bean",
		WrongAnswer2:  "Monk's Cafe",
		WrongAnswer3:  "The Grind",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "pop-003",
		PackID:        PackOfficePop,
		Category:      CategoryTVFilm,
		QuestionText:  "Which movie franchise features Woody and Buzz Lightyear?",
		CorrectAnswer: "Toy Story",
		WrongAnswer1:  "Cars",
		WrongAnswer2:  "Shrek",
		WrongAnswer3:  "Despicable Me",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "pop-004",
		PackID:        PackOfficePop,
		Category:      CategoryMusicCulture,
		QuestionText:  `Which sport is often called "the beautiful game"?`,
		CorrectAnswer: "Soccer",
		WrongAnswer1:  "Basketball",
		WrongAnswer2:  "Tennis",
		WrongAnswer3:  "Baseball",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "pop-005",
		PackID:        PackOfficePop,
		Category:      CategoryMusicCulture,
		QuestionText:  `Which artist released the song "Shake It Off"?`,
		CorrectAnswer: "Taylor Swift",
		WrongAnswer1:  "Katy Perry",
		WrongAnswer2:  "Ariana Grande",
		WrongAnswer3:  "Dua Lipa",
		MinRating:     contentrating.Work,
	},
	{
		ID:            "pop-006",
		PackID:        PackOfficePop,
		Category:      CategoryTVFilm,
		QuestionText:  "What is the name of the school in the Harry Potter series?",
		CorrectAnswer: "Hogwarts",
		WrongAnswer1:  "Beauxbatons",
		WrongAnswer2:  "Durmstrang",
		WrongAnswer3:  "Ilvermorny",
		MinRating:     contentrating.Work,
	},

	// Quick Brain Boost
	{
		ID:            "quick-001",
		PackID:        PackQuickBrain,
		Category:      CategoryScienceMath,
		QuestionText:  "Which planet is known as the Red Planet?",
		CorrectAnswer: "Mars",
		WrongAnswer1:  "Venus",
		WrongAnswer2:  "Jupiter",
		WrongAnswer3:  "Mercury",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "quick-002",
		PackID:        PackQuickBrain,
		Category:      CategoryEverydayFacts,
		QuestionText:  "What is the largest ocean on Earth?",
		CorrectAnswer: "Pacific Ocean",
		WrongAnswer1:  "Atlantic Ocean",
		WrongAnswer2:  "Indian Ocean",
		WrongAnswer3:  "Arctic Ocean",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "quick-003",
		PackID:        PackQuickBrain,
		Category:      CategoryScienceMath,
		QuestionText:  "How many sides does a hexagon have?",
		CorrectAnswer: "6",
		WrongAnswer1:  "5",
		WrongAnswer2:  "7",
		WrongAnswer3:  "8",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "quick-004",
		PackID:        PackQuickBrain,
		Category:      CategoryScienceMath,
		QuestionText:  "What is H2O commonly called?",
		CorrectAnswer: "Water",
		WrongAnswer1:  "Hydrogen Peroxide",
		WrongAnswer2:  "Salt",
		WrongAnswer3:  "Ozone",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "quick-005",
		PackID:        PackQuickBrain,
		Category:      CategoryEverydayFacts,
		QuestionText:  "What is the capital city of Japan?",
		CorrectAnswer: "Tokyo",
		WrongAnswer1:  "Kyoto",
		WrongAnswer2:  "Osaka",
		WrongAnswer3:  "Seoul",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "quick-006",
		PackID:        PackQuickBrain,
		Category:      CategoryScienceMath,
		QuestionText:  "Which instrument has 88 keys on a standard model?",
		CorrectAnswer: "Piano",
		WrongAnswer1:  "Guitar",
		WrongAnswer2:  "Violin",
		WrongAnswer3:  "Saxophone",
		MinRating:     contentrating.Kids,
	},

	// World Snapshot
	{
		ID:            "world-001",
		PackID:        PackWorldSnapshot,
		Category:      CategoryHistory,
		QuestionText:  "Who was the first person to walk on the moon?",
		CorrectAnswer: "Neil Armstrong",
		WrongAnswer1:  "Buzz Aldrin",
		WrongAnswer2:  "Yuri Gagarin",
		WrongAnswer3:  "John Glenn",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "world-002",
		PackID:        PackWorldSnapshot,
		Category:      CategoryHistory,
		QuestionText:  "What year did the Berlin Wall fall?",
		CorrectAnswer: "1989",
		WrongAnswer1:  "1987",
		WrongAnswer2:  "1991",
		WrongAnswer3:  "1979",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "world-003",
		PackID:        PackWorldSnapshot,
		Category:      CategoryHistory,
		QuestionText:  "In what year did World War II end?",
		CorrectAnswer: "1945",
		WrongAnswer1:  "1944",
		WrongAnswer2:  "1946",
		WrongAnswer3:  "1939",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "world-004",
		PackID:        PackWorldSnapshot,
		Category:      CategoryGeography,
		QuestionText:  "What is the smallest country in the world by area?",
		CorrectAnswer: "Vatican City",
		WrongAnswer1:  "Monaco",
		WrongAnswer2:  "San Marino",
		WrongAnswer3:  "Liechtenstein",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "world-005",
		PackID:        PackWorldSnapshot,
		Category:      CategoryHistory,
		QuestionText:  "Which country gifted the Statue of Liberty to the United States?",
		CorrectAnswer: "France",
		WrongAnswer1:  "United Kingdom",
		WrongAnswer2:  "Spain",
		WrongAnswer3:  "Italy",
		MinRating:     contentrating.Kids,
	},
	{
		ID:            "world-006",
		PackID:        PackWorldSnapshot,
		Category:      CategoryGeography,
		QuestionText:  "On which continent is the Sahara Desert located?",
		CorrectAnswer: "Africa",
		WrongAnswer1:  "Asia",
		WrongAnswer2:  "Australia",
		WrongAnswer3:  "South America",
		MinRating:     contentrating.Kids,
	},
}

func GetTemplateByID(id string) *Template {
	for i := range AllTemplates {
		if AllTemplates[i].ID == id {
			return &AllTemplates[i]
		}
	}
	return nil
}

func GetPackByID(id string) *Pack {
	for i := range AllPacks {
		if AllPacks[i].ID == id {
			return &AllPacks[i]
		}
	}
	return nil
}

// GetAvailableTemplates returns templates that are unused and allowed for the lobby audience.
func GetAvailableTemplates(usedIDs []string, lobbyContentRating int16) []Template {
	usedSet := make(map[string]bool, len(usedIDs))
	for _, id := range usedIDs {
		usedSet[id] = true
	}

	available := make([]Template, 0, len(AllTemplates))
	for _, t := range AllTemplates {
		if usedSet[t.ID] {
			continue
		}
		if t.MinRating > lobbyContentRating {
			continue
		}
		available = append(available, t)
	}

	return available
}

// BuildPackSections returns templates grouped first by pack, then by category.
func BuildPackSections(usedIDs []string, lobbyContentRating int16) []PackSection {
	available := GetAvailableTemplates(usedIDs, lobbyContentRating)
	templatesByPack := make(map[string][]Template)
	for _, t := range available {
		templatesByPack[t.PackID] = append(templatesByPack[t.PackID], t)
	}

	sections := make([]PackSection, 0, len(AllPacks))
	for _, pack := range AllPacks {
		if pack.MinRating > lobbyContentRating {
			continue
		}
		templates := templatesByPack[pack.ID]
		if len(templates) == 0 {
			continue
		}

		grouped := GroupByCategory(templates)
		sections = append(sections, PackSection{
			Pack:                pack,
			Categories:          orderedCategories(grouped),
			TemplatesByCategory: grouped,
			TemplateCount:       len(templates),
		})
	}

	return sections
}

// GroupByCategory groups templates by their category.
func GroupByCategory(templates []Template) map[string][]Template {
	grouped := make(map[string][]Template)
	for _, t := range templates {
		grouped[t.Category] = append(grouped[t.Category], t)
	}
	return grouped
}

func orderedCategories(grouped map[string][]Template) []string {
	result := make([]string, 0, len(grouped))
	seen := make(map[string]bool, len(grouped))

	for _, category := range categoryOrder {
		if _, ok := grouped[category]; ok {
			result = append(result, category)
			seen[category] = true
		}
	}

	extra := make([]string, 0)
	for category := range grouped {
		if !seen[category] {
			extra = append(extra, category)
		}
	}
	sort.Strings(extra)
	result = append(result, extra...)

	return result
}
