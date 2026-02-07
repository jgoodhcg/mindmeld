package questions

// Template represents a pre-made question template that players can use
type Template struct {
	ID            string
	Category      string
	QuestionText  string
	CorrectAnswer string
	WrongAnswer1  string
	WrongAnswer2  string
	WrongAnswer3  string
}

// Categories for organizing templates
const (
	CategoryPopCulture  = "Pop Culture"
	CategoryHistory     = "History"
	CategoryScience     = "Science & Nature"
	CategoryGeography   = "Geography"
	CategoryPersonal    = "Personal"
	CategoryPreferences = "About Me"
)

// AllTemplates contains all available question templates
var AllTemplates = []Template{
	// Pop Culture
	{
		ID:            "pop-001",
		Category:      CategoryPopCulture,
		QuestionText:  "Which movie won the Academy Award for Best Picture in 2020?",
		CorrectAnswer: "Parasite",
		WrongAnswer1:  "1917",
		WrongAnswer2:  "Joker",
		WrongAnswer3:  "Once Upon a Time in Hollywood",
	},
	{
		ID:            "pop-002",
		Category:      CategoryPopCulture,
		QuestionText:  "What is the highest-grossing film of all time (not adjusted for inflation)?",
		CorrectAnswer: "Avatar",
		WrongAnswer1:  "Avengers: Endgame",
		WrongAnswer2:  "Titanic",
		WrongAnswer3:  "Star Wars: The Force Awakens",
	},
	{
		ID:            "pop-003",
		Category:      CategoryPopCulture,
		QuestionText:  "Which band released the album 'Abbey Road'?",
		CorrectAnswer: "The Beatles",
		WrongAnswer1:  "The Rolling Stones",
		WrongAnswer2:  "Led Zeppelin",
		WrongAnswer3:  "Pink Floyd",
	},
	{
		ID:            "pop-004",
		Category:      CategoryPopCulture,
		QuestionText:  "What TV show features a chemistry teacher turned drug manufacturer?",
		CorrectAnswer: "Breaking Bad",
		WrongAnswer1:  "Better Call Saul",
		WrongAnswer2:  "Ozark",
		WrongAnswer3:  "The Wire",
	},
	{
		ID:            "pop-005",
		Category:      CategoryPopCulture,
		QuestionText:  "Who played Iron Man in the Marvel Cinematic Universe?",
		CorrectAnswer: "Robert Downey Jr.",
		WrongAnswer1:  "Chris Evans",
		WrongAnswer2:  "Chris Hemsworth",
		WrongAnswer3:  "Mark Ruffalo",
	},
	{
		ID:            "pop-006",
		Category:      CategoryPopCulture,
		QuestionText:  "Which artist released the album 'Thriller'?",
		CorrectAnswer: "Michael Jackson",
		WrongAnswer1:  "Prince",
		WrongAnswer2:  "Whitney Houston",
		WrongAnswer3:  "Madonna",
	},

	// History
	{
		ID:            "hist-001",
		Category:      CategoryHistory,
		QuestionText:  "In what year did World War II end?",
		CorrectAnswer: "1945",
		WrongAnswer1:  "1944",
		WrongAnswer2:  "1946",
		WrongAnswer3:  "1943",
	},
	{
		ID:            "hist-002",
		Category:      CategoryHistory,
		QuestionText:  "Who was the first President of the United States?",
		CorrectAnswer: "George Washington",
		WrongAnswer1:  "Thomas Jefferson",
		WrongAnswer2:  "John Adams",
		WrongAnswer3:  "Benjamin Franklin",
	},
	{
		ID:            "hist-003",
		Category:      CategoryHistory,
		QuestionText:  "Which ancient wonder was located in Alexandria, Egypt?",
		CorrectAnswer: "The Lighthouse (Pharos)",
		WrongAnswer1:  "The Hanging Gardens",
		WrongAnswer2:  "The Colossus",
		WrongAnswer3:  "The Great Pyramid",
	},
	{
		ID:            "hist-004",
		Category:      CategoryHistory,
		QuestionText:  "What year did the Berlin Wall fall?",
		CorrectAnswer: "1989",
		WrongAnswer1:  "1991",
		WrongAnswer2:  "1987",
		WrongAnswer3:  "1990",
	},
	{
		ID:            "hist-005",
		Category:      CategoryHistory,
		QuestionText:  "Who was the first person to walk on the moon?",
		CorrectAnswer: "Neil Armstrong",
		WrongAnswer1:  "Buzz Aldrin",
		WrongAnswer2:  "John Glenn",
		WrongAnswer3:  "Yuri Gagarin",
	},
	{
		ID:            "hist-006",
		Category:      CategoryHistory,
		QuestionText:  "Which empire was ruled by Julius Caesar?",
		CorrectAnswer: "Roman Empire",
		WrongAnswer1:  "Greek Empire",
		WrongAnswer2:  "Persian Empire",
		WrongAnswer3:  "Ottoman Empire",
	},

	// Science & Nature
	{
		ID:            "sci-001",
		Category:      CategoryScience,
		QuestionText:  "What is the chemical symbol for gold?",
		CorrectAnswer: "Au",
		WrongAnswer1:  "Ag",
		WrongAnswer2:  "Go",
		WrongAnswer3:  "Gd",
	},
	{
		ID:            "sci-002",
		Category:      CategoryScience,
		QuestionText:  "What planet is known as the Red Planet?",
		CorrectAnswer: "Mars",
		WrongAnswer1:  "Venus",
		WrongAnswer2:  "Jupiter",
		WrongAnswer3:  "Mercury",
	},
	{
		ID:            "sci-003",
		Category:      CategoryScience,
		QuestionText:  "What is the largest organ in the human body?",
		CorrectAnswer: "Skin",
		WrongAnswer1:  "Liver",
		WrongAnswer2:  "Brain",
		WrongAnswer3:  "Heart",
	},
	{
		ID:            "sci-004",
		Category:      CategoryScience,
		QuestionText:  "How many bones are in the adult human body?",
		CorrectAnswer: "206",
		WrongAnswer1:  "186",
		WrongAnswer2:  "226",
		WrongAnswer3:  "256",
	},
	{
		ID:            "sci-005",
		Category:      CategoryScience,
		QuestionText:  "What gas do plants absorb from the atmosphere?",
		CorrectAnswer: "Carbon dioxide",
		WrongAnswer1:  "Oxygen",
		WrongAnswer2:  "Nitrogen",
		WrongAnswer3:  "Hydrogen",
	},
	{
		ID:            "sci-006",
		Category:      CategoryScience,
		QuestionText:  "What is the speed of light in a vacuum (approximately)?",
		CorrectAnswer: "300,000 km/s",
		WrongAnswer1:  "150,000 km/s",
		WrongAnswer2:  "500,000 km/s",
		WrongAnswer3:  "1,000,000 km/s",
	},

	// Geography
	{
		ID:            "geo-001",
		Category:      CategoryGeography,
		QuestionText:  "What is the capital of Australia?",
		CorrectAnswer: "Canberra",
		WrongAnswer1:  "Sydney",
		WrongAnswer2:  "Melbourne",
		WrongAnswer3:  "Brisbane",
	},
	{
		ID:            "geo-002",
		Category:      CategoryGeography,
		QuestionText:  "Which country has the most time zones?",
		CorrectAnswer: "France",
		WrongAnswer1:  "Russia",
		WrongAnswer2:  "United States",
		WrongAnswer3:  "China",
	},
	{
		ID:            "geo-003",
		Category:      CategoryGeography,
		QuestionText:  "What is the longest river in the world?",
		CorrectAnswer: "Nile",
		WrongAnswer1:  "Amazon",
		WrongAnswer2:  "Yangtze",
		WrongAnswer3:  "Mississippi",
	},
	{
		ID:            "geo-004",
		Category:      CategoryGeography,
		QuestionText:  "Which country is home to the Great Barrier Reef?",
		CorrectAnswer: "Australia",
		WrongAnswer1:  "Indonesia",
		WrongAnswer2:  "Philippines",
		WrongAnswer3:  "New Zealand",
	},
	{
		ID:            "geo-005",
		Category:      CategoryGeography,
		QuestionText:  "What is the smallest country in the world?",
		CorrectAnswer: "Vatican City",
		WrongAnswer1:  "Monaco",
		WrongAnswer2:  "San Marino",
		WrongAnswer3:  "Liechtenstein",
	},
	{
		ID:            "geo-006",
		Category:      CategoryGeography,
		QuestionText:  "On which continent is the Sahara Desert located?",
		CorrectAnswer: "Africa",
		WrongAnswer1:  "Asia",
		WrongAnswer2:  "Middle East",
		WrongAnswer3:  "Australia",
	},

	// Personal / Get to Know You - these have placeholder answers for players to customize
	{
		ID:            "personal-001",
		Category:      CategoryPersonal,
		QuestionText:  "What is [my name]'s favorite movie of all time?",
		CorrectAnswer: "[Your favorite movie]",
		WrongAnswer1:  "[Wrong option 1]",
		WrongAnswer2:  "[Wrong option 2]",
		WrongAnswer3:  "[Wrong option 3]",
	},
	{
		ID:            "personal-002",
		Category:      CategoryPersonal,
		QuestionText:  "What city was [my name] born in?",
		CorrectAnswer: "[Your birth city]",
		WrongAnswer1:  "[Wrong city 1]",
		WrongAnswer2:  "[Wrong city 2]",
		WrongAnswer3:  "[Wrong city 3]",
	},
	{
		ID:            "personal-003",
		Category:      CategoryPersonal,
		QuestionText:  "How many siblings does [my name] have?",
		CorrectAnswer: "[Your answer]",
		WrongAnswer1:  "[Wrong number 1]",
		WrongAnswer2:  "[Wrong number 2]",
		WrongAnswer3:  "[Wrong number 3]",
	},
	{
		ID:            "personal-004",
		Category:      CategoryPersonal,
		QuestionText:  "What was [my name]'s first job?",
		CorrectAnswer: "[Your first job]",
		WrongAnswer1:  "[Wrong job 1]",
		WrongAnswer2:  "[Wrong job 2]",
		WrongAnswer3:  "[Wrong job 3]",
	},
	{
		ID:            "personal-005",
		Category:      CategoryPersonal,
		QuestionText:  "What is [my name]'s dream vacation destination?",
		CorrectAnswer: "[Your dream destination]",
		WrongAnswer1:  "[Wrong destination 1]",
		WrongAnswer2:  "[Wrong destination 2]",
		WrongAnswer3:  "[Wrong destination 3]",
	},
	{
		ID:            "personal-006",
		Category:      CategoryPersonal,
		QuestionText:  "What is [my name]'s hidden talent?",
		CorrectAnswer: "[Your hidden talent]",
		WrongAnswer1:  "[Wrong talent 1]",
		WrongAnswer2:  "[Wrong talent 2]",
		WrongAnswer3:  "[Wrong talent 3]",
	},

	// About Me
	{
		ID:            "pref-001",
		Category:      CategoryPreferences,
		QuestionText:  "What is [my name]'s favorite cuisine?",
		CorrectAnswer: "[Your favorite cuisine]",
		WrongAnswer1:  "[Wrong cuisine 1]",
		WrongAnswer2:  "[Wrong cuisine 2]",
		WrongAnswer3:  "[Wrong cuisine 3]",
	},
	{
		ID:            "pref-002",
		Category:      CategoryPreferences,
		QuestionText:  "What is [my name]'s go-to comfort food?",
		CorrectAnswer: "[Your comfort food]",
		WrongAnswer1:  "[Wrong food 1]",
		WrongAnswer2:  "[Wrong food 2]",
		WrongAnswer3:  "[Wrong food 3]",
	},
	{
		ID:            "pref-003",
		Category:      CategoryPreferences,
		QuestionText:  "What music genre does [my name] listen to most?",
		CorrectAnswer: "[Your favorite genre]",
		WrongAnswer1:  "[Wrong genre 1]",
		WrongAnswer2:  "[Wrong genre 2]",
		WrongAnswer3:  "[Wrong genre 3]",
	},
	{
		ID:            "pref-004",
		Category:      CategoryPreferences,
		QuestionText:  "What is [my name]'s favorite season?",
		CorrectAnswer: "[Your favorite season]",
		WrongAnswer1:  "[Wrong season 1]",
		WrongAnswer2:  "[Wrong season 2]",
		WrongAnswer3:  "[Wrong season 3]",
	},
	{
		ID:            "pref-005",
		Category:      CategoryPreferences,
		QuestionText:  "Is [my name] a morning person or night owl?",
		CorrectAnswer: "[Morning person/Night owl]",
		WrongAnswer1:  "[The opposite]",
		WrongAnswer2:  "Neither",
		WrongAnswer3:  "Both equally",
	},
	{
		ID:            "pref-006",
		Category:      CategoryPreferences,
		QuestionText:  "What is [my name]'s favorite holiday?",
		CorrectAnswer: "[Your favorite holiday]",
		WrongAnswer1:  "[Wrong holiday 1]",
		WrongAnswer2:  "[Wrong holiday 2]",
		WrongAnswer3:  "[Wrong holiday 3]",
	},
}

// GetTemplateByID returns a template by its ID, or nil if not found
func GetTemplateByID(id string) *Template {
	for i := range AllTemplates {
		if AllTemplates[i].ID == id {
			return &AllTemplates[i]
		}
	}
	return nil
}

// GetAvailableTemplates returns templates that haven't been used (not in usedIDs)
func GetAvailableTemplates(usedIDs []string) []Template {
	usedSet := make(map[string]bool)
	for _, id := range usedIDs {
		usedSet[id] = true
	}

	available := make([]Template, 0)
	for _, t := range AllTemplates {
		if !usedSet[t.ID] {
			available = append(available, t)
		}
	}
	return available
}

// GetCategories returns all unique categories
func GetCategories() []string {
	return []string{
		CategoryPreferences,
		CategoryPersonal,
		CategoryPopCulture,
		CategoryHistory,
		CategoryScience,
		CategoryGeography,
	}
}

// GroupByCategory groups templates by their category
func GroupByCategory(templates []Template) map[string][]Template {
	grouped := make(map[string][]Template)
	for _, t := range templates {
		grouped[t.Category] = append(grouped[t.Category], t)
	}
	return grouped
}
