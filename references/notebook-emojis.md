```js

const notebookIcons = [
  // ============================================================================
  // GENERIC BOOKS & NOTEBOOKS
  // ============================================================================
  { emoji: "📚", keywords: ["general", "study", "books", "reading", "library", "research"] },
  { emoji: "📖", keywords: ["english", "literature", "reading", "writing", "language arts", "book"] },
  { emoji: "📕", keywords: ["general", "book", "notebook", "red", "study", "textbook"] },
  { emoji: "📗", keywords: ["general", "book", "notebook", "green", "study", "textbook"] },
  { emoji: "📘", keywords: ["general", "book", "notebook", "blue", "study", "textbook"] },
  { emoji: "📙", keywords: ["general", "book", "notebook", "orange", "study", "textbook"] },
  { emoji: "📓", keywords: ["notebook", "journal", "notes", "diary", "general"] },
  { emoji: "📔", keywords: ["notebook", "journal", "decorative", "diary", "general"] },
  { emoji: "📒", keywords: ["ledger", "notebook", "accounting", "notes", "yellow"] },
  { emoji: "🗒️", keywords: ["notepad", "spiral", "notes", "memo", "list", "general"] },
  { emoji: "📝", keywords: ["notes", "memo", "writing", "pencil", "general", "homework"] },
  { emoji: "🔖", keywords: ["bookmark", "reading", "reference", "save", "tag"] },
  { emoji: "📑", keywords: ["tabs", "bookmarks", "reference", "organization", "documents"] },
  { emoji: "🗂️", keywords: ["organization", "notes", "reference", "archive", "filing", "folders"] },
  { emoji: "📋", keywords: ["clipboard", "management", "project", "checklist", "planning", "organization"] },
  
  // ============================================================================
  // SCIENCES
  // ============================================================================
  { emoji: "🧬", keywords: ["biology", "genetics", "dna", "life science", "molecular", "biotech"] },
  { emoji: "⚛️", keywords: ["physics", "atom", "nuclear", "quantum", "science"] },
  { emoji: "🔬", keywords: ["chemistry", "lab", "science", "research", "experiment", "microscope"] },
  { emoji: "🧪", keywords: ["chemistry", "lab", "experiment", "test tube", "science"] },
  { emoji: "🦠", keywords: ["microbiology", "bacteria", "virus", "cells", "immunology"] },
  { emoji: "🧫", keywords: ["biology", "petri dish", "lab", "culture", "microbiology"] },
  { emoji: "🌱", keywords: ["botany", "plants", "ecology", "environmental", "agriculture", "growth"] },
  { emoji: "🌿", keywords: ["botany", "plants", "herbs", "nature", "ecology", "green"] },
  { emoji: "🍃", keywords: ["ecology", "nature", "environmental", "plants", "sustainability"] },
  { emoji: "🌸", keywords: ["botany", "flowers", "plants", "japanese", "spring", "nature"] },
  { emoji: "🌳", keywords: ["forestry", "ecology", "trees", "environmental", "nature"] },
  { emoji: "🍄", keywords: ["mycology", "fungi", "biology", "nature", "foraging"] },
  { emoji: "🌍", keywords: ["geography", "earth science", "environmental", "geology", "climate", "world"] },
  { emoji: "🌎", keywords: ["geography", "americas", "earth", "world", "global studies"] },
  { emoji: "🌏", keywords: ["geography", "asia", "pacific", "international relations", "global studies"] },
  { emoji: "🌋", keywords: ["geology", "earth science", "volcanology", "rocks", "minerals"] },
  { emoji: "⛰️", keywords: ["geology", "geography", "mountains", "earth science", "hiking"] },
  { emoji: "🏔️", keywords: ["geology", "geography", "mountains", "alpine", "glaciology"] },
  { emoji: "💎", keywords: ["geology", "mineralogy", "gems", "crystals", "earth science"] },
  { emoji: "🪨", keywords: ["geology", "rocks", "earth science", "minerals", "petrology"] },
  { emoji: "🌊", keywords: ["oceanography", "marine biology", "waves", "water", "hydrology"] },
  { emoji: "🐚", keywords: ["marine biology", "ocean", "shells", "beach", "zoology"] },
  { emoji: "🐋", keywords: ["marine biology", "whales", "ocean", "zoology", "mammals"] },
  { emoji: "🦈", keywords: ["marine biology", "sharks", "ocean", "zoology", "fish"] },
  { emoji: "🐠", keywords: ["marine biology", "fish", "aquarium", "ichthyology", "ocean"] },
  { emoji: "🔭", keywords: ["astronomy", "space", "stars", "astrophysics", "telescope"] },
  { emoji: "🚀", keywords: ["aerospace", "space", "astronomy", "rockets", "physics"] },
  { emoji: "🪐", keywords: ["astronomy", "planets", "space", "saturn", "solar system"] },
  { emoji: "🌙", keywords: ["astronomy", "moon", "lunar", "space", "night"] },
  { emoji: "⭐", keywords: ["astronomy", "stars", "space", "general", "favorites"] },
  { emoji: "☀️", keywords: ["astronomy", "solar", "sun", "energy", "physics"] },
  { emoji: "🌡️", keywords: ["thermodynamics", "temperature", "weather", "climate", "physics"] },
  { emoji: "⚡", keywords: ["electricity", "physics", "energy", "electrical engineering", "power"] },
  { emoji: "🧲", keywords: ["physics", "magnetism", "electromagnetics", "science"] },
  { emoji: "⏳", keywords: ["history", "time", "physics", "hourglass", "ancient"] },
  { emoji: "🕰️", keywords: ["history", "time", "horology", "clocks", "antique"] },
  
  // ============================================================================
  // MATH & LOGIC
  // ============================================================================
  { emoji: "🧮", keywords: ["math", "mathematics", "calculation", "arithmetic", "accounting", "abacus"] },
  { emoji: "📐", keywords: ["geometry", "math", "architecture", "trigonometry", "angles", "drafting"] },
  { emoji: "📊", keywords: ["statistics", "data", "analytics", "graphs", "business", "charts"] },
  { emoji: "📈", keywords: ["economics", "finance", "growth", "statistics", "business", "trends"] },
  { emoji: "📉", keywords: ["economics", "finance", "decline", "statistics", "analysis"] },
  { emoji: "🔢", keywords: ["math", "numbers", "algebra", "calculus", "arithmetic"] },
  { emoji: "♾️", keywords: ["math", "calculus", "infinity", "limits", "theory"] },
  { emoji: "➕", keywords: ["math", "addition", "arithmetic", "basic math", "elementary"] },
  { emoji: "🧩", keywords: ["logic", "puzzles", "problem solving", "games", "critical thinking"] },
  { emoji: "🎯", keywords: ["goals", "targets", "focus", "precision", "planning"] },
  
  // ============================================================================
  // TECHNOLOGY & ENGINEERING
  // ============================================================================
  { emoji: "💻", keywords: ["computer science", "programming", "coding", "software", "tech"] },
  { emoji: "🖥️", keywords: ["computer science", "desktop", "technology", "it", "software"] },
  { emoji: "⌨️", keywords: ["computer science", "typing", "programming", "keyboard", "tech"] },
  { emoji: "🖱️", keywords: ["computer science", "technology", "ui", "ux", "interface"] },
  { emoji: "⚙️", keywords: ["engineering", "mechanical", "machines", "systems", "mechanics", "settings"] },
  { emoji: "🔧", keywords: ["engineering", "tools", "mechanical", "repair", "technical"] },
  { emoji: "🔩", keywords: ["engineering", "hardware", "mechanical", "construction", "fasteners"] },
  { emoji: "🛠️", keywords: ["engineering", "tools", "workshop", "technical", "repair", "diy"] },
  { emoji: "🤖", keywords: ["robotics", "ai", "artificial intelligence", "machine learning", "automation"] },
  { emoji: "🔌", keywords: ["electrical engineering", "electronics", "circuits", "hardware", "power"] },
  { emoji: "💡", keywords: ["ideas", "innovation", "electrical", "invention", "brainstorm", "project"] },
  { emoji: "🔋", keywords: ["electrical engineering", "energy", "batteries", "power", "electronics"] },
  { emoji: "📡", keywords: ["telecommunications", "signals", "radio", "networking", "broadcast"] },
  { emoji: "🌐", keywords: ["web development", "internet", "networking", "languages", "global"] },
  { emoji: "📱", keywords: ["mobile development", "apps", "technology", "ux", "design"] },
  { emoji: "🔐", keywords: ["cybersecurity", "security", "cryptography", "infosec", "privacy"] },
  { emoji: "🔑", keywords: ["security", "cryptography", "access", "keys", "authentication"] },
  { emoji: "🏗️", keywords: ["construction", "civil engineering", "building", "architecture", "structural"] },
  { emoji: "🧱", keywords: ["construction", "building", "materials", "masonry", "civil engineering"] },
  { emoji: "🏠", keywords: ["architecture", "housing", "real estate", "design", "home"] },
  { emoji: "🏛️", keywords: ["architecture", "history", "philosophy", "politics", "government", "classics"] },
  { emoji: "✈️", keywords: ["aerospace", "aviation", "travel", "flight", "transportation"] },
  { emoji: "🛩️", keywords: ["aviation", "aerospace", "pilot", "flight", "aircraft"] },
  { emoji: "🚁", keywords: ["aviation", "aerospace", "helicopter", "flight", "engineering"] },
  { emoji: "🚗", keywords: ["automotive", "transportation", "cars", "engineering", "mechanics"] },
  { emoji: "🚂", keywords: ["transportation", "trains", "railways", "logistics", "engineering"] },
  { emoji: "🚢", keywords: ["maritime", "naval", "shipping", "boats", "oceanography"] },
  { emoji: "⚓", keywords: ["maritime", "naval", "nautical", "sailing", "ocean"] },
  
  // ============================================================================
  // MEDICINE & HEALTH
  // ============================================================================
  { emoji: "🩺", keywords: ["medicine", "health", "doctor", "medical", "clinical", "diagnosis"] },
  { emoji: "🫀", keywords: ["anatomy", "cardiology", "heart", "physiology", "medicine"] },
  { emoji: "🫁", keywords: ["anatomy", "pulmonology", "lungs", "respiratory", "medicine"] },
  { emoji: "🧠", keywords: ["neuroscience", "psychology", "brain", "cognitive", "mental health"] },
  { emoji: "👁️", keywords: ["ophthalmology", "optometry", "vision", "eyes", "anatomy"] },
  { emoji: "👂", keywords: ["audiology", "hearing", "ear", "anatomy", "ent"] },
  { emoji: "🦷", keywords: ["dentistry", "dental", "teeth", "oral health", "orthodontics"] },
  { emoji: "🦴", keywords: ["anatomy", "orthopedics", "skeleton", "physiology", "osteology"] },
  { emoji: "💪", keywords: ["anatomy", "muscles", "fitness", "physiology", "kinesiology"] },
  { emoji: "🩸", keywords: ["hematology", "blood", "medicine", "phlebotomy", "lab"] },
  { emoji: "💉", keywords: ["medicine", "vaccines", "injections", "nursing", "immunology"] },
  { emoji: "💊", keywords: ["pharmacology", "medicine", "drugs", "pharmacy", "pharmaceutical"] },
  { emoji: "🏥", keywords: ["nursing", "healthcare", "hospital", "clinical", "medical"] },
  { emoji: "🧘", keywords: ["wellness", "yoga", "meditation", "mindfulness", "mental health"] },
  { emoji: "🥗", keywords: ["nutrition", "dietetics", "health", "food", "wellness"] },
  { emoji: "🏋️", keywords: ["fitness", "exercise", "kinesiology", "sports medicine", "training"] },
  { emoji: "😷", keywords: ["epidemiology", "public health", "medicine", "healthcare", "infection"] },
  { emoji: "🧬", keywords: ["genetics", "genomics", "dna", "heredity", "molecular biology"] },
  
  // ============================================================================
  // HUMANITIES & SOCIAL SCIENCES
  // ============================================================================
  { emoji: "📜", keywords: ["history", "ancient", "classics", "documents", "archives", "scrolls"] },
  { emoji: "🏺", keywords: ["archaeology", "ancient", "history", "anthropology", "artifacts"] },
  { emoji: "🗿", keywords: ["archaeology", "anthropology", "ancient", "statues", "history"] },
  { emoji: "⚱️", keywords: ["archaeology", "ancient", "history", "artifacts", "ceramics"] },
  { emoji: "🦕", keywords: ["paleontology", "dinosaurs", "fossils", "prehistoric", "geology"] },
  { emoji: "🦖", keywords: ["paleontology", "dinosaurs", "fossils", "prehistoric", "natural history"] },
  { emoji: "✍️", keywords: ["writing", "composition", "creative writing", "essay", "journalism"] },
  { emoji: "🖋️", keywords: ["writing", "literature", "calligraphy", "poetry", "fountain pen"] },
  { emoji: "✒️", keywords: ["writing", "literature", "poetry", "classics", "penmanship"] },
  { emoji: "🎭", keywords: ["theater", "drama", "acting", "performing arts", "literature"] },
  { emoji: "📰", keywords: ["journalism", "news", "media", "current events", "press"] },
  { emoji: "🗞️", keywords: ["journalism", "news", "newspaper", "media", "current events"] },
  { emoji: "⚖️", keywords: ["law", "ethics", "justice", "legal", "philosophy", "politics"] },
  { emoji: "🗳️", keywords: ["political science", "government", "civics", "elections", "democracy"] },
  { emoji: "🏴", keywords: ["political science", "flags", "nations", "government", "activism"] },
  { emoji: "💬", keywords: ["linguistics", "communication", "speech", "language", "debate"] },
  { emoji: "🗣️", keywords: ["speech", "communication", "linguistics", "public speaking", "debate"] },
  { emoji: "👥", keywords: ["sociology", "social work", "community", "groups", "society"] },
  { emoji: "🤝", keywords: ["social work", "cooperation", "partnerships", "diplomacy", "relations"] },
  { emoji: "❤️", keywords: ["social work", "health", "relationships", "psychology", "care"] },
  { emoji: "👶", keywords: ["child development", "pediatrics", "early childhood", "education", "family"] },
  { emoji: "👨‍👩‍👧‍👦", keywords: ["family studies", "sociology", "psychology", "relationships", "social work"] },
  { emoji: "🧓", keywords: ["gerontology", "aging", "elder care", "sociology", "health"] },
  
  // ============================================================================
  // RELIGION & PHILOSOPHY
  // ============================================================================
  { emoji: "🙏", keywords: ["religion", "spirituality", "prayer", "theology", "faith"] },
  { emoji: "☯️", keywords: ["philosophy", "taoism", "eastern", "balance", "religion"] },
  { emoji: "🕉️", keywords: ["hinduism", "religion", "spirituality", "sanskrit", "yoga"] },
  { emoji: "☸️", keywords: ["buddhism", "religion", "spirituality", "dharma", "meditation"] },
  { emoji: "✡️", keywords: ["judaism", "religion", "hebrew", "theology", "spirituality"] },
  { emoji: "✝️", keywords: ["christianity", "religion", "theology", "bible", "spirituality"] },
  { emoji: "☪️", keywords: ["islam", "religion", "arabic", "theology", "spirituality"] },
  { emoji: "🔮", keywords: ["mysticism", "occult", "divination", "spirituality", "metaphysics"] },
  { emoji: "🪬", keywords: ["protection", "culture", "spirituality", "symbolism", "middle eastern"] },
  
  // ============================================================================
  // BUSINESS & ECONOMICS
  // ============================================================================
  { emoji: "💰", keywords: ["economics", "finance", "money", "business", "banking", "accounting"] },
  { emoji: "💵", keywords: ["finance", "money", "economics", "cash", "business"] },
  { emoji: "💳", keywords: ["finance", "payments", "banking", "credit", "business"] },
  { emoji: "💼", keywords: ["business", "management", "corporate", "mba", "entrepreneurship"] },
  { emoji: "🏦", keywords: ["finance", "banking", "economics", "investment", "accounting"] },
  { emoji: "🏢", keywords: ["business", "corporate", "office", "management", "real estate"] },
  { emoji: "🏪", keywords: ["retail", "business", "commerce", "marketing", "entrepreneurship"] },
  { emoji: "🛒", keywords: ["retail", "commerce", "marketing", "consumer", "business"] },
  { emoji: "📦", keywords: ["logistics", "supply chain", "shipping", "inventory", "business"] },
  { emoji: "🏷️", keywords: ["marketing", "branding", "retail", "pricing", "business"] },
  
  // ============================================================================
  // ARTS & DESIGN
  // ============================================================================
  { emoji: "🎨", keywords: ["art", "design", "painting", "creative", "visual arts", "drawing"] },
  { emoji: "🖼️", keywords: ["art", "painting", "gallery", "museum", "visual arts", "frames"] },
  { emoji: "🖌️", keywords: ["art", "painting", "brush", "creative", "illustration"] },
  { emoji: "✏️", keywords: ["drawing", "sketching", "design", "illustration", "art", "writing"] },
  { emoji: "🖍️", keywords: ["art", "drawing", "coloring", "illustration", "elementary"] },
  { emoji: "🎭", keywords: ["theater", "drama", "acting", "masks", "performing arts"] },
  { emoji: "🎬", keywords: ["film", "cinema", "video", "media", "production", "movies"] },
  { emoji: "🎥", keywords: ["film", "video", "cinema", "production", "recording"] },
  { emoji: "📷", keywords: ["photography", "visual", "media", "film", "art", "camera"] },
  { emoji: "📸", keywords: ["photography", "camera", "visual", "media", "snapshot"] },
  { emoji: "🎞️", keywords: ["film", "cinema", "photography", "media", "animation"] },
  { emoji: "🎵", keywords: ["music", "theory", "composition", "audio", "sound"] },
  { emoji: "🎶", keywords: ["music", "songs", "melody", "composition", "audio"] },
  { emoji: "🎼", keywords: ["music", "sheet music", "composition", "theory", "notation"] },
  { emoji: "🎹", keywords: ["piano", "music", "instrument", "theory", "composition", "keyboard"] },
  { emoji: "🎸", keywords: ["guitar", "music", "instrument", "rock", "strings"] },
  { emoji: "🎺", keywords: ["trumpet", "music", "brass", "instrument", "jazz"] },
  { emoji: "🎻", keywords: ["violin", "music", "strings", "orchestra", "classical"] },
  { emoji: "🥁", keywords: ["drums", "music", "percussion", "instrument", "rhythm"] },
  { emoji: "🎤", keywords: ["vocals", "singing", "music", "performance", "karaoke"] },
  { emoji: "🎧", keywords: ["audio", "music", "sound", "production", "listening"] },
  { emoji: "🎙️", keywords: ["podcast", "radio", "broadcasting", "audio", "journalism"] },
  { emoji: "📻", keywords: ["radio", "broadcasting", "audio", "media", "communications"] },
  { emoji: "📺", keywords: ["television", "media", "broadcast", "film", "communications"] },
  { emoji: "💃", keywords: ["dance", "performing arts", "movement", "choreography", "ballet"] },
  { emoji: "🩰", keywords: ["ballet", "dance", "performing arts", "choreography", "movement"] },
  { emoji: "👗", keywords: ["fashion", "design", "textiles", "clothing", "style"] },
  { emoji: "🧵", keywords: ["textiles", "sewing", "fashion", "crafts", "fabrics"] },
  { emoji: "✂️", keywords: ["crafts", "cutting", "sewing", "design", "editing"] },
  { emoji: "🪡", keywords: ["sewing", "textiles", "crafts", "embroidery", "fashion"] },
  { emoji: "🧶", keywords: ["knitting", "crafts", "textiles", "fiber arts", "crochet"] },
  
  // ============================================================================
  // ANIMALS & ZOOLOGY
  // ============================================================================
  { emoji: "🐾", keywords: ["zoology", "animals", "veterinary", "pets", "wildlife"] },
  { emoji: "🦋", keywords: ["entomology", "insects", "butterflies", "biology", "nature"] },
  { emoji: "🐝", keywords: ["entomology", "bees", "insects", "ecology", "agriculture"] },
  { emoji: "🐜", keywords: ["entomology", "ants", "insects", "biology", "myrmecology"] },
  { emoji: "🐛", keywords: ["entomology", "insects", "caterpillars", "biology", "larvae"] },
  { emoji: "🦅", keywords: ["ornithology", "birds", "eagles", "wildlife", "zoology"] },
  { emoji: "🦉", keywords: ["ornithology", "birds", "owls", "wildlife", "nocturnal"] },
  { emoji: "🐦", keywords: ["ornithology", "birds", "wildlife", "zoology", "nature"] },
  { emoji: "🐺", keywords: ["zoology", "wolves", "wildlife", "mammals", "ecology"] },
  { emoji: "🦁", keywords: ["zoology", "lions", "wildlife", "mammals", "safari"] },
  { emoji: "🐻", keywords: ["zoology", "bears", "wildlife", "mammals", "nature"] },
  { emoji: "🐴", keywords: ["equine", "horses", "veterinary", "animals", "riding"] },
  { emoji: "🐶", keywords: ["veterinary", "dogs", "pets", "animals", "canine"] },
  { emoji: "🐱", keywords: ["veterinary", "cats", "pets", "animals", "feline"] },
  { emoji: "🦎", keywords: ["herpetology", "reptiles", "lizards", "zoology", "biology"] },
  { emoji: "🐍", keywords: ["herpetology", "snakes", "reptiles", "zoology", "biology"] },
  { emoji: "🐸", keywords: ["herpetology", "frogs", "amphibians", "zoology", "biology"] },
  
  // ============================================================================
  // AGRICULTURE & ENVIRONMENT
  // ============================================================================
  { emoji: "🌾", keywords: ["agriculture", "farming", "crops", "grain", "food science"] },
  { emoji: "🚜", keywords: ["agriculture", "farming", "tractors", "machinery", "rural"] },
  { emoji: "🌻", keywords: ["agriculture", "botany", "flowers", "sunflowers", "nature"] },
  { emoji: "🌽", keywords: ["agriculture", "crops", "corn", "farming", "food"] },
  { emoji: "🍇", keywords: ["viticulture", "wine", "grapes", "agriculture", "horticulture"] },
  { emoji: "🍎", keywords: ["agriculture", "fruit", "horticulture", "nutrition", "food"] },
  { emoji: "♻️", keywords: ["environmental", "sustainability", "recycling", "ecology", "green"] },
  { emoji: "🌲", keywords: ["forestry", "trees", "nature", "environmental", "ecology"] },
  { emoji: "🏕️", keywords: ["outdoor education", "camping", "nature", "recreation", "wilderness"] },
  { emoji: "🥾", keywords: ["outdoor education", "hiking", "recreation", "nature", "adventure"] },
  
  // ============================================================================
  // FOOD & CULINARY
  // ============================================================================
  { emoji: "🍳", keywords: ["culinary", "cooking", "food science", "nutrition", "hospitality"] },
  { emoji: "👨‍🍳", keywords: ["culinary", "chef", "cooking", "hospitality", "gastronomy"] },
  { emoji: "🍰", keywords: ["baking", "pastry", "desserts", "culinary", "patisserie"] },
  { emoji: "🥘", keywords: ["culinary", "cooking", "cuisine", "gastronomy", "food"] },
  { emoji: "🍷", keywords: ["wine", "sommelier", "viticulture", "hospitality", "beverage"] },
  { emoji: "☕", keywords: ["coffee", "barista", "beverage", "hospitality", "cafe"] },
  { emoji: "🍵", keywords: ["tea", "beverage", "culture", "hospitality", "ceremony"] },
  { emoji: "🧁", keywords: ["baking", "desserts", "pastry", "culinary", "decorating"] },
  
  // ============================================================================
  // SPORTS & PHYSICAL EDUCATION
  // ============================================================================
  { emoji: "🏃", keywords: ["physical education", "sports", "kinesiology", "fitness", "athletics", "running"] },
  { emoji: "⚽", keywords: ["soccer", "football", "sports", "athletics", "team sports"] },
  { emoji: "🏀", keywords: ["basketball", "sports", "athletics", "team sports", "nba"] },
  { emoji: "🏈", keywords: ["football", "american football", "sports", "athletics", "nfl"] },
  { emoji: "⚾", keywords: ["baseball", "sports", "athletics", "team sports", "mlb"] },
  { emoji: "🎾", keywords: ["tennis", "sports", "athletics", "racket", "individual sports"] },
  { emoji: "🏐", keywords: ["volleyball", "sports", "athletics", "team sports", "beach"] },
  { emoji: "🏊", keywords: ["swimming", "sports", "athletics", "aquatics", "fitness"] },
  { emoji: "🚴", keywords: ["cycling", "sports", "fitness", "athletics", "biking"] },
  { emoji: "🧗", keywords: ["climbing", "sports", "adventure", "fitness", "outdoor"] },
  { emoji: "🥋", keywords: ["martial arts", "karate", "judo", "sports", "self defense"] },
  { emoji: "🤺", keywords: ["fencing", "sports", "athletics", "swordsmanship", "olympic"] },
  { emoji: "♟️", keywords: ["chess", "strategy", "games", "logic", "competition"] },
  { emoji: "🎲", keywords: ["games", "probability", "statistics", "board games", "chance"] },
  { emoji: "🎮", keywords: ["gaming", "video games", "game design", "esports", "interactive"] },
  { emoji: "🎯", keywords: ["archery", "darts", "precision", "sports", "targeting"] },
  { emoji: "🏆", keywords: ["competition", "awards", "sports", "achievement", "championship"] },
  { emoji: "🥇", keywords: ["competition", "achievement", "sports", "first place", "excellence"] },
  
  // ============================================================================
  // EDUCATION & ACADEMIC
  // ============================================================================
  { emoji: "🎓", keywords: ["graduation", "academic", "degree", "university", "college", "exam"] },
  { emoji: "🏫", keywords: ["school", "education", "learning", "teaching", "academic"] },
  { emoji: "👩‍🏫", keywords: ["teaching", "education", "instructor", "professor", "lecture"] },
  { emoji: "👨‍🎓", keywords: ["student", "education", "university", "college", "learning"] },
  { emoji: "🎒", keywords: ["school", "student", "education", "backpack", "supplies"] },
  { emoji: "📏", keywords: ["measurement", "math", "geometry", "ruler", "precision"] },
  { emoji: "🔎", keywords: ["research", "investigation", "analysis", "study", "search"] },
  { emoji: "🔍", keywords: ["research", "investigation", "analysis", "study", "detail"] },
  
  // ============================================================================
  // TRAVEL & GEOGRAPHY
  // ============================================================================
  { emoji: "🗺️", keywords: ["geography", "maps", "travel", "cartography", "exploration"] },
  { emoji: "🧭", keywords: ["navigation", "geography", "compass", "orientation", "exploration"] },
  { emoji: "🗽", keywords: ["american studies", "history", "travel", "landmarks", "usa"] },
  { emoji: "🗼", keywords: ["architecture", "landmarks", "travel", "paris", "structures"] },
  { emoji: "🏰", keywords: ["history", "castles", "medieval", "architecture", "european"] },
  { emoji: "⛩️", keywords: ["japanese", "culture", "religion", "architecture", "asian studies"] },
  { emoji: "🕌", keywords: ["islamic studies", "architecture", "religion", "culture", "mosque"] },
  { emoji: "⛪", keywords: ["religion", "christianity", "architecture", "church", "theology"] },
  { emoji: "🛕", keywords: ["hinduism", "temple", "religion", "architecture", "indian"] },
  
  // ============================================================================
  // MYTHOLOGY & FANTASY
  // ============================================================================
  { emoji: "🐉", keywords: ["mythology", "dragons", "fantasy", "legends", "folklore"] },
  { emoji: "🦄", keywords: ["mythology", "unicorn", "fantasy", "legends", "fairy tales"] },
  { emoji: "🧙", keywords: ["fantasy", "magic", "wizard", "mythology", "folklore"] },
  { emoji: "🧚", keywords: ["fantasy", "fairy", "mythology", "folklore", "fairy tales"] },
  { emoji: "🧜", keywords: ["mythology", "mermaids", "fantasy", "folklore", "ocean"] },
  { emoji: "🪄", keywords: ["magic", "fantasy", "illusion", "tricks", "wonder"] },
  { emoji: "✨", keywords: ["magic", "sparkle", "special", "creativity", "highlights"] },
  
  // ============================================================================
  // MILITARY & DEFENSE
  // ============================================================================
  { emoji: "🎖️", keywords: ["military", "medals", "defense", "honors", "service"] },
  { emoji: "🛡️", keywords: ["defense", "security", "protection", "military", "shields"] },
  { emoji: "⚔️", keywords: ["military history", "combat", "swords", "warfare", "medieval"] },
  { emoji: "🪖", keywords: ["military", "army", "defense", "helmet", "soldier"] },
  
  // ============================================================================
  // MISCELLANEOUS
  // ============================================================================
  { emoji: "🔥", keywords: ["trending", "hot topics", "popular", "fire", "urgent"] },
  { emoji: "💫", keywords: ["highlights", "special", "important", "dizzy", "stars"] },
  { emoji: "🌈", keywords: ["diversity", "lgbtq", "pride", "colors", "inclusivity"] },
  { emoji: "🎪", keywords: ["circus", "entertainment", "performing arts", "events", "shows"] },
  { emoji: "🎠", keywords: ["amusement", "recreation", "entertainment", "fun", "parks"] },
  { emoji: "🎡", keywords: ["amusement", "recreation", "entertainment", "ferris wheel", "parks"] },
  { emoji: "🎢", keywords: ["amusement", "recreation", "entertainment", "roller coaster", "physics"] },
  { emoji: "🌺", keywords: ["hawaiian", "tropical", "flowers", "culture", "nature"] },
  { emoji: "🐕‍🦺", keywords: ["service animals", "assistance", "disability studies", "therapy", "support"] },
  { emoji: "🤟", keywords: ["sign language", "deaf studies", "communication", "asl", "accessibility"] },
  { emoji: "♿", keywords: ["disability studies", "accessibility", "inclusion", "accommodation", "ada"] },
];
```

