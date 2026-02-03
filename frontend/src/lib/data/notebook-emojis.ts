// Notebook emoji data for the emoji selector component
// Categories are organized for easy browsing when selecting a notebook icon

export interface NotebookEmoji {
	emoji: string;
	keywords: string[];
	category: string;
}

export const notebookEmojis: NotebookEmoji[] = [
	// Generic Books & Notebooks
	{ emoji: '📚', keywords: ['general', 'study', 'books', 'reading', 'library', 'research'], category: 'Books & Notebooks' },
	{ emoji: '📖', keywords: ['english', 'literature', 'reading', 'writing', 'language arts', 'book'], category: 'Books & Notebooks' },
	{ emoji: '📕', keywords: ['general', 'book', 'notebook', 'red', 'study', 'textbook'], category: 'Books & Notebooks' },
	{ emoji: '📗', keywords: ['general', 'book', 'notebook', 'green', 'study', 'textbook'], category: 'Books & Notebooks' },
	{ emoji: '📘', keywords: ['general', 'book', 'notebook', 'blue', 'study', 'textbook'], category: 'Books & Notebooks' },
	{ emoji: '📙', keywords: ['general', 'book', 'notebook', 'orange', 'study', 'textbook'], category: 'Books & Notebooks' },
	{ emoji: '📓', keywords: ['notebook', 'journal', 'notes', 'diary', 'general'], category: 'Books & Notebooks' },
	{ emoji: '📔', keywords: ['notebook', 'journal', 'decorative', 'diary', 'general'], category: 'Books & Notebooks' },
	{ emoji: '📒', keywords: ['ledger', 'notebook', 'accounting', 'notes', 'yellow'], category: 'Books & Notebooks' },
	{ emoji: '🗒️', keywords: ['notepad', 'spiral', 'notes', 'memo', 'list', 'general'], category: 'Books & Notebooks' },
	{ emoji: '📝', keywords: ['notes', 'memo', 'writing', 'pencil', 'general', 'homework'], category: 'Books & Notebooks' },
	{ emoji: '🔖', keywords: ['bookmark', 'reading', 'reference', 'save', 'tag'], category: 'Books & Notebooks' },
	{ emoji: '📑', keywords: ['tabs', 'bookmarks', 'reference', 'organization', 'documents'], category: 'Books & Notebooks' },
	{ emoji: '🗂️', keywords: ['organization', 'notes', 'reference', 'archive', 'filing', 'folders'], category: 'Books & Notebooks' },
	{ emoji: '📋', keywords: ['clipboard', 'management', 'project', 'checklist', 'planning', 'organization'], category: 'Books & Notebooks' },

	// Sciences
	{ emoji: '🧬', keywords: ['biology', 'genetics', 'dna', 'life science', 'molecular', 'biotech'], category: 'Sciences' },
	{ emoji: '⚛️', keywords: ['physics', 'atom', 'nuclear', 'quantum', 'science'], category: 'Sciences' },
	{ emoji: '🔬', keywords: ['chemistry', 'lab', 'science', 'research', 'experiment', 'microscope'], category: 'Sciences' },
	{ emoji: '🧪', keywords: ['chemistry', 'lab', 'experiment', 'test tube', 'science'], category: 'Sciences' },
	{ emoji: '🦠', keywords: ['microbiology', 'bacteria', 'virus', 'cells', 'immunology'], category: 'Sciences' },
	{ emoji: '🧫', keywords: ['biology', 'petri dish', 'lab', 'culture', 'microbiology'], category: 'Sciences' },
	{ emoji: '🌱', keywords: ['botany', 'plants', 'ecology', 'environmental', 'agriculture', 'growth'], category: 'Sciences' },
	{ emoji: '🌿', keywords: ['botany', 'plants', 'herbs', 'nature', 'ecology', 'green'], category: 'Sciences' },
	{ emoji: '🍃', keywords: ['ecology', 'nature', 'environmental', 'plants', 'sustainability'], category: 'Sciences' },
	{ emoji: '🌸', keywords: ['botany', 'flowers', 'plants', 'japanese', 'spring', 'nature'], category: 'Sciences' },
	{ emoji: '🌳', keywords: ['forestry', 'ecology', 'trees', 'environmental', 'nature'], category: 'Sciences' },
	{ emoji: '🍄', keywords: ['mycology', 'fungi', 'biology', 'nature', 'foraging'], category: 'Sciences' },
	{ emoji: '🌍', keywords: ['geography', 'earth science', 'environmental', 'geology', 'climate', 'world'], category: 'Sciences' },
	{ emoji: '🌎', keywords: ['geography', 'americas', 'earth', 'world', 'global studies'], category: 'Sciences' },
	{ emoji: '🌏', keywords: ['geography', 'asia', 'pacific', 'international relations', 'global studies'], category: 'Sciences' },
	{ emoji: '🌋', keywords: ['geology', 'earth science', 'volcanology', 'rocks', 'minerals'], category: 'Sciences' },
	{ emoji: '⛰️', keywords: ['geology', 'geography', 'mountains', 'earth science', 'hiking'], category: 'Sciences' },
	{ emoji: '🏔️', keywords: ['geology', 'geography', 'mountains', 'alpine', 'glaciology'], category: 'Sciences' },
	{ emoji: '💎', keywords: ['geology', 'mineralogy', 'gems', 'crystals', 'earth science'], category: 'Sciences' },
	{ emoji: '🪨', keywords: ['geology', 'rocks', 'earth science', 'minerals', 'petrology'], category: 'Sciences' },
	{ emoji: '🌊', keywords: ['oceanography', 'marine biology', 'waves', 'water', 'hydrology'], category: 'Sciences' },
	{ emoji: '🐚', keywords: ['marine biology', 'ocean', 'shells', 'beach', 'zoology'], category: 'Sciences' },
	{ emoji: '🐋', keywords: ['marine biology', 'whales', 'ocean', 'zoology', 'mammals'], category: 'Sciences' },
	{ emoji: '🦈', keywords: ['marine biology', 'sharks', 'ocean', 'zoology', 'fish'], category: 'Sciences' },
	{ emoji: '🐠', keywords: ['marine biology', 'fish', 'aquarium', 'ichthyology', 'ocean'], category: 'Sciences' },
	{ emoji: '🔭', keywords: ['astronomy', 'space', 'stars', 'astrophysics', 'telescope'], category: 'Sciences' },
	{ emoji: '🚀', keywords: ['aerospace', 'space', 'astronomy', 'rockets', 'physics'], category: 'Sciences' },
	{ emoji: '🪐', keywords: ['astronomy', 'planets', 'space', 'saturn', 'solar system'], category: 'Sciences' },
	{ emoji: '🌙', keywords: ['astronomy', 'moon', 'lunar', 'space', 'night'], category: 'Sciences' },
	{ emoji: '⭐', keywords: ['astronomy', 'stars', 'space', 'general', 'favorites'], category: 'Sciences' },
	{ emoji: '☀️', keywords: ['astronomy', 'solar', 'sun', 'energy', 'physics'], category: 'Sciences' },
	{ emoji: '🌡️', keywords: ['thermodynamics', 'temperature', 'weather', 'climate', 'physics'], category: 'Sciences' },
	{ emoji: '⚡', keywords: ['electricity', 'physics', 'energy', 'electrical engineering', 'power'], category: 'Sciences' },
	{ emoji: '🧲', keywords: ['physics', 'magnetism', 'electromagnetics', 'science'], category: 'Sciences' },
	{ emoji: '⏳', keywords: ['history', 'time', 'physics', 'hourglass', 'ancient'], category: 'Sciences' },
	{ emoji: '🕰️', keywords: ['history', 'time', 'horology', 'clocks', 'antique'], category: 'Sciences' },

	// Math & Logic
	{ emoji: '🧮', keywords: ['math', 'mathematics', 'calculation', 'arithmetic', 'accounting', 'abacus'], category: 'Math & Logic' },
	{ emoji: '📐', keywords: ['geometry', 'math', 'architecture', 'trigonometry', 'angles', 'drafting'], category: 'Math & Logic' },
	{ emoji: '📊', keywords: ['statistics', 'data', 'analytics', 'graphs', 'business', 'charts'], category: 'Math & Logic' },
	{ emoji: '📈', keywords: ['economics', 'finance', 'growth', 'statistics', 'business', 'trends'], category: 'Math & Logic' },
	{ emoji: '📉', keywords: ['economics', 'finance', 'decline', 'statistics', 'analysis'], category: 'Math & Logic' },
	{ emoji: '🔢', keywords: ['math', 'numbers', 'algebra', 'calculus', 'arithmetic'], category: 'Math & Logic' },
	{ emoji: '♾️', keywords: ['math', 'calculus', 'infinity', 'limits', 'theory'], category: 'Math & Logic' },
	{ emoji: '➕', keywords: ['math', 'addition', 'arithmetic', 'basic math', 'elementary'], category: 'Math & Logic' },
	{ emoji: '🧩', keywords: ['logic', 'puzzles', 'problem solving', 'games', 'critical thinking'], category: 'Math & Logic' },
	{ emoji: '🎯', keywords: ['goals', 'targets', 'focus', 'precision', 'planning'], category: 'Math & Logic' },

	// Technology & Engineering
	{ emoji: '💻', keywords: ['computer science', 'programming', 'coding', 'software', 'tech'], category: 'Technology & Engineering' },
	{ emoji: '🖥️', keywords: ['computer science', 'desktop', 'technology', 'it', 'software'], category: 'Technology & Engineering' },
	{ emoji: '⌨️', keywords: ['computer science', 'typing', 'programming', 'keyboard', 'tech'], category: 'Technology & Engineering' },
	{ emoji: '🖱️', keywords: ['computer science', 'technology', 'ui', 'ux', 'interface'], category: 'Technology & Engineering' },
	{ emoji: '⚙️', keywords: ['engineering', 'mechanical', 'machines', 'systems', 'mechanics', 'settings'], category: 'Technology & Engineering' },
	{ emoji: '🔧', keywords: ['engineering', 'tools', 'mechanical', 'repair', 'technical'], category: 'Technology & Engineering' },
	{ emoji: '🔩', keywords: ['engineering', 'hardware', 'mechanical', 'construction', 'fasteners'], category: 'Technology & Engineering' },
	{ emoji: '🛠️', keywords: ['engineering', 'tools', 'workshop', 'technical', 'repair', 'diy'], category: 'Technology & Engineering' },
	{ emoji: '🤖', keywords: ['robotics', 'ai', 'artificial intelligence', 'machine learning', 'automation'], category: 'Technology & Engineering' },
	{ emoji: '🔌', keywords: ['electrical engineering', 'electronics', 'circuits', 'hardware', 'power'], category: 'Technology & Engineering' },
	{ emoji: '💡', keywords: ['ideas', 'innovation', 'electrical', 'invention', 'brainstorm', 'project'], category: 'Technology & Engineering' },
	{ emoji: '🔋', keywords: ['electrical engineering', 'energy', 'batteries', 'power', 'electronics'], category: 'Technology & Engineering' },
	{ emoji: '📡', keywords: ['telecommunications', 'signals', 'radio', 'networking', 'broadcast'], category: 'Technology & Engineering' },
	{ emoji: '🌐', keywords: ['web development', 'internet', 'networking', 'languages', 'global'], category: 'Technology & Engineering' },
	{ emoji: '📱', keywords: ['mobile development', 'apps', 'technology', 'ux', 'design'], category: 'Technology & Engineering' },
	{ emoji: '🔐', keywords: ['cybersecurity', 'security', 'cryptography', 'infosec', 'privacy'], category: 'Technology & Engineering' },
	{ emoji: '🔑', keywords: ['security', 'cryptography', 'access', 'keys', 'authentication'], category: 'Technology & Engineering' },
	{ emoji: '🏗️', keywords: ['construction', 'civil engineering', 'building', 'architecture', 'structural'], category: 'Technology & Engineering' },
	{ emoji: '🧱', keywords: ['construction', 'building', 'materials', 'masonry', 'civil engineering'], category: 'Technology & Engineering' },
	{ emoji: '🏠', keywords: ['architecture', 'housing', 'real estate', 'design', 'home'], category: 'Technology & Engineering' },
	{ emoji: '🏛️', keywords: ['architecture', 'history', 'philosophy', 'politics', 'government', 'classics'], category: 'Technology & Engineering' },
	{ emoji: '✈️', keywords: ['aerospace', 'aviation', 'travel', 'flight', 'transportation'], category: 'Technology & Engineering' },
	{ emoji: '🛩️', keywords: ['aviation', 'aerospace', 'pilot', 'flight', 'aircraft'], category: 'Technology & Engineering' },
	{ emoji: '🚁', keywords: ['aviation', 'aerospace', 'helicopter', 'flight', 'engineering'], category: 'Technology & Engineering' },
	{ emoji: '🚗', keywords: ['automotive', 'transportation', 'cars', 'engineering', 'mechanics'], category: 'Technology & Engineering' },
	{ emoji: '🚂', keywords: ['transportation', 'trains', 'railways', 'logistics', 'engineering'], category: 'Technology & Engineering' },
	{ emoji: '🚢', keywords: ['maritime', 'naval', 'shipping', 'boats', 'oceanography'], category: 'Technology & Engineering' },
	{ emoji: '⚓', keywords: ['maritime', 'naval', 'nautical', 'sailing', 'ocean'], category: 'Technology & Engineering' },

	// Medicine & Health
	{ emoji: '🩺', keywords: ['medicine', 'health', 'doctor', 'medical', 'clinical', 'diagnosis'], category: 'Medicine & Health' },
	{ emoji: '🫀', keywords: ['anatomy', 'cardiology', 'heart', 'physiology', 'medicine'], category: 'Medicine & Health' },
	{ emoji: '🫁', keywords: ['anatomy', 'pulmonology', 'lungs', 'respiratory', 'medicine'], category: 'Medicine & Health' },
	{ emoji: '🧠', keywords: ['neuroscience', 'psychology', 'brain', 'cognitive', 'mental health'], category: 'Medicine & Health' },
	{ emoji: '👁️', keywords: ['ophthalmology', 'optometry', 'vision', 'eyes', 'anatomy'], category: 'Medicine & Health' },
	{ emoji: '👂', keywords: ['audiology', 'hearing', 'ear', 'anatomy', 'ent'], category: 'Medicine & Health' },
	{ emoji: '🦷', keywords: ['dentistry', 'dental', 'teeth', 'oral health', 'orthodontics'], category: 'Medicine & Health' },
	{ emoji: '🦴', keywords: ['anatomy', 'orthopedics', 'skeleton', 'physiology', 'osteology'], category: 'Medicine & Health' },
	{ emoji: '💪', keywords: ['anatomy', 'muscles', 'fitness', 'physiology', 'kinesiology'], category: 'Medicine & Health' },
	{ emoji: '🩸', keywords: ['hematology', 'blood', 'medicine', 'phlebotomy', 'lab'], category: 'Medicine & Health' },
	{ emoji: '💉', keywords: ['medicine', 'vaccines', 'injections', 'nursing', 'immunology'], category: 'Medicine & Health' },
	{ emoji: '💊', keywords: ['pharmacology', 'medicine', 'drugs', 'pharmacy', 'pharmaceutical'], category: 'Medicine & Health' },
	{ emoji: '🏥', keywords: ['nursing', 'healthcare', 'hospital', 'clinical', 'medical'], category: 'Medicine & Health' },
	{ emoji: '🧘', keywords: ['wellness', 'yoga', 'meditation', 'mindfulness', 'mental health'], category: 'Medicine & Health' },
	{ emoji: '🥗', keywords: ['nutrition', 'dietetics', 'health', 'food', 'wellness'], category: 'Medicine & Health' },
	{ emoji: '🏋️', keywords: ['fitness', 'exercise', 'kinesiology', 'sports medicine', 'training'], category: 'Medicine & Health' },
	{ emoji: '😷', keywords: ['epidemiology', 'public health', 'medicine', 'healthcare', 'infection'], category: 'Medicine & Health' },

	// Humanities & Social Sciences
	{ emoji: '📜', keywords: ['history', 'ancient', 'classics', 'documents', 'archives', 'scrolls'], category: 'Humanities & Social Sciences' },
	{ emoji: '🏺', keywords: ['archaeology', 'ancient', 'history', 'anthropology', 'artifacts'], category: 'Humanities & Social Sciences' },
	{ emoji: '🗿', keywords: ['archaeology', 'anthropology', 'ancient', 'statues', 'history'], category: 'Humanities & Social Sciences' },
	{ emoji: '⚱️', keywords: ['archaeology', 'ancient', 'history', 'artifacts', 'ceramics'], category: 'Humanities & Social Sciences' },
	{ emoji: '🦕', keywords: ['paleontology', 'dinosaurs', 'fossils', 'prehistoric', 'geology'], category: 'Humanities & Social Sciences' },
	{ emoji: '🦖', keywords: ['paleontology', 'dinosaurs', 'fossils', 'prehistoric', 'natural history'], category: 'Humanities & Social Sciences' },
	{ emoji: '✍️', keywords: ['writing', 'composition', 'creative writing', 'essay', 'journalism'], category: 'Humanities & Social Sciences' },
	{ emoji: '🖋️', keywords: ['writing', 'literature', 'calligraphy', 'poetry', 'fountain pen'], category: 'Humanities & Social Sciences' },
	{ emoji: '✒️', keywords: ['writing', 'literature', 'poetry', 'classics', 'penmanship'], category: 'Humanities & Social Sciences' },
	{ emoji: '🎭', keywords: ['theater', 'drama', 'acting', 'performing arts', 'literature'], category: 'Humanities & Social Sciences' },
	{ emoji: '📰', keywords: ['journalism', 'news', 'media', 'current events', 'press'], category: 'Humanities & Social Sciences' },
	{ emoji: '🗞️', keywords: ['journalism', 'news', 'newspaper', 'media', 'current events'], category: 'Humanities & Social Sciences' },
	{ emoji: '⚖️', keywords: ['law', 'ethics', 'justice', 'legal', 'philosophy', 'politics'], category: 'Humanities & Social Sciences' },
	{ emoji: '🗳️', keywords: ['political science', 'government', 'civics', 'elections', 'democracy'], category: 'Humanities & Social Sciences' },
	{ emoji: '🏴', keywords: ['political science', 'flags', 'nations', 'government', 'activism'], category: 'Humanities & Social Sciences' },
	{ emoji: '💬', keywords: ['linguistics', 'communication', 'speech', 'language', 'debate'], category: 'Humanities & Social Sciences' },
	{ emoji: '🗣️', keywords: ['speech', 'communication', 'linguistics', 'public speaking', 'debate'], category: 'Humanities & Social Sciences' },
	{ emoji: '👥', keywords: ['sociology', 'social work', 'community', 'groups', 'society'], category: 'Humanities & Social Sciences' },
	{ emoji: '🤝', keywords: ['social work', 'cooperation', 'partnerships', 'diplomacy', 'relations'], category: 'Humanities & Social Sciences' },
	{ emoji: '❤️', keywords: ['social work', 'health', 'relationships', 'psychology', 'care'], category: 'Humanities & Social Sciences' },
	{ emoji: '👶', keywords: ['child development', 'pediatrics', 'early childhood', 'education', 'family'], category: 'Humanities & Social Sciences' },
	{ emoji: '👨‍👩‍👧‍👦', keywords: ['family studies', 'sociology', 'psychology', 'relationships', 'social work'], category: 'Humanities & Social Sciences' },
	{ emoji: '🧓', keywords: ['gerontology', 'aging', 'elder care', 'sociology', 'health'], category: 'Humanities & Social Sciences' },

	// Religion & Philosophy
	{ emoji: '🙏', keywords: ['religion', 'spirituality', 'prayer', 'theology', 'faith'], category: 'Religion & Philosophy' },
	{ emoji: '☯️', keywords: ['philosophy', 'taoism', 'eastern', 'balance', 'religion'], category: 'Religion & Philosophy' },
	{ emoji: '🕉️', keywords: ['hinduism', 'religion', 'spirituality', 'sanskrit', 'yoga'], category: 'Religion & Philosophy' },
	{ emoji: '☸️', keywords: ['buddhism', 'religion', 'spirituality', 'dharma', 'meditation'], category: 'Religion & Philosophy' },
	{ emoji: '✡️', keywords: ['judaism', 'religion', 'hebrew', 'theology', 'spirituality'], category: 'Religion & Philosophy' },
	{ emoji: '✝️', keywords: ['christianity', 'religion', 'theology', 'bible', 'spirituality'], category: 'Religion & Philosophy' },
	{ emoji: '☪️', keywords: ['islam', 'religion', 'arabic', 'theology', 'spirituality'], category: 'Religion & Philosophy' },
	{ emoji: '🔮', keywords: ['mysticism', 'occult', 'divination', 'spirituality', 'metaphysics'], category: 'Religion & Philosophy' },
	{ emoji: '🪬', keywords: ['protection', 'culture', 'spirituality', 'symbolism', 'middle eastern'], category: 'Religion & Philosophy' },

	// Business & Economics
	{ emoji: '💰', keywords: ['economics', 'finance', 'money', 'business', 'banking', 'accounting'], category: 'Business & Economics' },
	{ emoji: '💵', keywords: ['finance', 'money', 'economics', 'cash', 'business'], category: 'Business & Economics' },
	{ emoji: '💳', keywords: ['finance', 'payments', 'banking', 'credit', 'business'], category: 'Business & Economics' },
	{ emoji: '💼', keywords: ['business', 'management', 'corporate', 'mba', 'entrepreneurship'], category: 'Business & Economics' },
	{ emoji: '🏦', keywords: ['finance', 'banking', 'economics', 'investment', 'accounting'], category: 'Business & Economics' },
	{ emoji: '🏢', keywords: ['business', 'corporate', 'office', 'management', 'real estate'], category: 'Business & Economics' },
	{ emoji: '🏪', keywords: ['retail', 'business', 'commerce', 'marketing', 'entrepreneurship'], category: 'Business & Economics' },
	{ emoji: '🛒', keywords: ['retail', 'commerce', 'marketing', 'consumer', 'business'], category: 'Business & Economics' },
	{ emoji: '📦', keywords: ['logistics', 'supply chain', 'shipping', 'inventory', 'business'], category: 'Business & Economics' },
	{ emoji: '🏷️', keywords: ['marketing', 'branding', 'retail', 'pricing', 'business'], category: 'Business & Economics' },

	// Arts & Design
	{ emoji: '🎨', keywords: ['art', 'design', 'painting', 'creative', 'visual arts', 'drawing'], category: 'Arts & Design' },
	{ emoji: '🖼️', keywords: ['art', 'painting', 'gallery', 'museum', 'visual arts', 'frames'], category: 'Arts & Design' },
	{ emoji: '🖌️', keywords: ['art', 'painting', 'brush', 'creative', 'illustration'], category: 'Arts & Design' },
	{ emoji: '✏️', keywords: ['drawing', 'sketching', 'design', 'illustration', 'art', 'writing'], category: 'Arts & Design' },
	{ emoji: '🖍️', keywords: ['art', 'drawing', 'coloring', 'illustration', 'elementary'], category: 'Arts & Design' },
	{ emoji: '🎬', keywords: ['film', 'cinema', 'video', 'media', 'production', 'movies'], category: 'Arts & Design' },
	{ emoji: '🎥', keywords: ['film', 'video', 'cinema', 'production', 'recording'], category: 'Arts & Design' },
	{ emoji: '📷', keywords: ['photography', 'visual', 'media', 'film', 'art', 'camera'], category: 'Arts & Design' },
	{ emoji: '📸', keywords: ['photography', 'camera', 'visual', 'media', 'snapshot'], category: 'Arts & Design' },
	{ emoji: '🎞️', keywords: ['film', 'cinema', 'photography', 'media', 'animation'], category: 'Arts & Design' },
	{ emoji: '🎵', keywords: ['music', 'theory', 'composition', 'audio', 'sound'], category: 'Arts & Design' },
	{ emoji: '🎶', keywords: ['music', 'songs', 'melody', 'composition', 'audio'], category: 'Arts & Design' },
	{ emoji: '🎼', keywords: ['music', 'sheet music', 'composition', 'theory', 'notation'], category: 'Arts & Design' },
	{ emoji: '🎹', keywords: ['piano', 'music', 'instrument', 'theory', 'composition', 'keyboard'], category: 'Arts & Design' },
	{ emoji: '🎸', keywords: ['guitar', 'music', 'instrument', 'rock', 'strings'], category: 'Arts & Design' },
	{ emoji: '🎺', keywords: ['trumpet', 'music', 'brass', 'instrument', 'jazz'], category: 'Arts & Design' },
	{ emoji: '🎻', keywords: ['violin', 'music', 'strings', 'orchestra', 'classical'], category: 'Arts & Design' },
	{ emoji: '🥁', keywords: ['drums', 'music', 'percussion', 'instrument', 'rhythm'], category: 'Arts & Design' },
	{ emoji: '🎤', keywords: ['vocals', 'singing', 'music', 'performance', 'karaoke'], category: 'Arts & Design' },
	{ emoji: '🎧', keywords: ['audio', 'music', 'sound', 'production', 'listening'], category: 'Arts & Design' },
	{ emoji: '🎙️', keywords: ['podcast', 'radio', 'broadcasting', 'audio', 'journalism'], category: 'Arts & Design' },
	{ emoji: '📻', keywords: ['radio', 'broadcasting', 'audio', 'media', 'communications'], category: 'Arts & Design' },
	{ emoji: '📺', keywords: ['television', 'media', 'broadcast', 'film', 'communications'], category: 'Arts & Design' },
	{ emoji: '💃', keywords: ['dance', 'performing arts', 'movement', 'choreography', 'ballet'], category: 'Arts & Design' },
	{ emoji: '🩰', keywords: ['ballet', 'dance', 'performing arts', 'choreography', 'movement'], category: 'Arts & Design' },
	{ emoji: '👗', keywords: ['fashion', 'design', 'textiles', 'clothing', 'style'], category: 'Arts & Design' },
	{ emoji: '🧵', keywords: ['textiles', 'sewing', 'fashion', 'crafts', 'fabrics'], category: 'Arts & Design' },
	{ emoji: '✂️', keywords: ['crafts', 'cutting', 'sewing', 'design', 'editing'], category: 'Arts & Design' },
	{ emoji: '🪡', keywords: ['sewing', 'textiles', 'crafts', 'embroidery', 'fashion'], category: 'Arts & Design' },
	{ emoji: '🧶', keywords: ['knitting', 'crafts', 'textiles', 'fiber arts', 'crochet'], category: 'Arts & Design' },

	// Animals & Zoology
	{ emoji: '🐾', keywords: ['zoology', 'animals', 'veterinary', 'pets', 'wildlife'], category: 'Animals & Zoology' },
	{ emoji: '🦋', keywords: ['entomology', 'insects', 'butterflies', 'biology', 'nature'], category: 'Animals & Zoology' },
	{ emoji: '🐝', keywords: ['entomology', 'bees', 'insects', 'ecology', 'agriculture'], category: 'Animals & Zoology' },
	{ emoji: '🐜', keywords: ['entomology', 'ants', 'insects', 'biology', 'myrmecology'], category: 'Animals & Zoology' },
	{ emoji: '🐛', keywords: ['entomology', 'insects', 'caterpillars', 'biology', 'larvae'], category: 'Animals & Zoology' },
	{ emoji: '🦅', keywords: ['ornithology', 'birds', 'eagles', 'wildlife', 'zoology'], category: 'Animals & Zoology' },
	{ emoji: '🦉', keywords: ['ornithology', 'birds', 'owls', 'wildlife', 'nocturnal'], category: 'Animals & Zoology' },
	{ emoji: '🐦', keywords: ['ornithology', 'birds', 'wildlife', 'zoology', 'nature'], category: 'Animals & Zoology' },
	{ emoji: '🐺', keywords: ['zoology', 'wolves', 'wildlife', 'mammals', 'ecology'], category: 'Animals & Zoology' },
	{ emoji: '🦁', keywords: ['zoology', 'lions', 'wildlife', 'mammals', 'safari'], category: 'Animals & Zoology' },
	{ emoji: '🐻', keywords: ['zoology', 'bears', 'wildlife', 'mammals', 'nature'], category: 'Animals & Zoology' },
	{ emoji: '🐴', keywords: ['equine', 'horses', 'veterinary', 'animals', 'riding'], category: 'Animals & Zoology' },
	{ emoji: '🐶', keywords: ['veterinary', 'dogs', 'pets', 'animals', 'canine'], category: 'Animals & Zoology' },
	{ emoji: '🐱', keywords: ['veterinary', 'cats', 'pets', 'animals', 'feline'], category: 'Animals & Zoology' },
	{ emoji: '🦎', keywords: ['herpetology', 'reptiles', 'lizards', 'zoology', 'biology'], category: 'Animals & Zoology' },
	{ emoji: '🐍', keywords: ['herpetology', 'snakes', 'reptiles', 'zoology', 'biology'], category: 'Animals & Zoology' },
	{ emoji: '🐸', keywords: ['herpetology', 'frogs', 'amphibians', 'zoology', 'biology'], category: 'Animals & Zoology' },

	// Agriculture & Environment
	{ emoji: '🌾', keywords: ['agriculture', 'farming', 'crops', 'grain', 'food science'], category: 'Agriculture & Environment' },
	{ emoji: '🚜', keywords: ['agriculture', 'farming', 'tractors', 'machinery', 'rural'], category: 'Agriculture & Environment' },
	{ emoji: '🌻', keywords: ['agriculture', 'botany', 'flowers', 'sunflowers', 'nature'], category: 'Agriculture & Environment' },
	{ emoji: '🌽', keywords: ['agriculture', 'crops', 'corn', 'farming', 'food'], category: 'Agriculture & Environment' },
	{ emoji: '🍇', keywords: ['viticulture', 'wine', 'grapes', 'agriculture', 'horticulture'], category: 'Agriculture & Environment' },
	{ emoji: '🍎', keywords: ['agriculture', 'fruit', 'horticulture', 'nutrition', 'food'], category: 'Agriculture & Environment' },
	{ emoji: '♻️', keywords: ['environmental', 'sustainability', 'recycling', 'ecology', 'green'], category: 'Agriculture & Environment' },
	{ emoji: '🌲', keywords: ['forestry', 'trees', 'nature', 'environmental', 'ecology'], category: 'Agriculture & Environment' },
	{ emoji: '🏕️', keywords: ['outdoor education', 'camping', 'nature', 'recreation', 'wilderness'], category: 'Agriculture & Environment' },
	{ emoji: '🥾', keywords: ['outdoor education', 'hiking', 'recreation', 'nature', 'adventure'], category: 'Agriculture & Environment' },

	// Food & Culinary
	{ emoji: '🍳', keywords: ['culinary', 'cooking', 'food science', 'nutrition', 'hospitality'], category: 'Food & Culinary' },
	{ emoji: '👨‍🍳', keywords: ['culinary', 'chef', 'cooking', 'hospitality', 'gastronomy'], category: 'Food & Culinary' },
	{ emoji: '🍰', keywords: ['baking', 'pastry', 'desserts', 'culinary', 'patisserie'], category: 'Food & Culinary' },
	{ emoji: '🥘', keywords: ['culinary', 'cooking', 'cuisine', 'gastronomy', 'food'], category: 'Food & Culinary' },
	{ emoji: '🍷', keywords: ['wine', 'sommelier', 'viticulture', 'hospitality', 'beverage'], category: 'Food & Culinary' },
	{ emoji: '☕', keywords: ['coffee', 'barista', 'beverage', 'hospitality', 'cafe'], category: 'Food & Culinary' },
	{ emoji: '🍵', keywords: ['tea', 'beverage', 'culture', 'hospitality', 'ceremony'], category: 'Food & Culinary' },
	{ emoji: '🧁', keywords: ['baking', 'desserts', 'pastry', 'culinary', 'decorating'], category: 'Food & Culinary' },

	// Sports & Physical Education
	{ emoji: '🏃', keywords: ['physical education', 'sports', 'kinesiology', 'fitness', 'athletics', 'running'], category: 'Sports & Fitness' },
	{ emoji: '⚽', keywords: ['soccer', 'football', 'sports', 'athletics', 'team sports'], category: 'Sports & Fitness' },
	{ emoji: '🏀', keywords: ['basketball', 'sports', 'athletics', 'team sports', 'nba'], category: 'Sports & Fitness' },
	{ emoji: '🏈', keywords: ['football', 'american football', 'sports', 'athletics', 'nfl'], category: 'Sports & Fitness' },
	{ emoji: '⚾', keywords: ['baseball', 'sports', 'athletics', 'team sports', 'mlb'], category: 'Sports & Fitness' },
	{ emoji: '🎾', keywords: ['tennis', 'sports', 'athletics', 'racket', 'individual sports'], category: 'Sports & Fitness' },
	{ emoji: '🏐', keywords: ['volleyball', 'sports', 'athletics', 'team sports', 'beach'], category: 'Sports & Fitness' },
	{ emoji: '🏊', keywords: ['swimming', 'sports', 'athletics', 'aquatics', 'fitness'], category: 'Sports & Fitness' },
	{ emoji: '🚴', keywords: ['cycling', 'sports', 'fitness', 'athletics', 'biking'], category: 'Sports & Fitness' },
	{ emoji: '🧗', keywords: ['climbing', 'sports', 'adventure', 'fitness', 'outdoor'], category: 'Sports & Fitness' },
	{ emoji: '🥋', keywords: ['martial arts', 'karate', 'judo', 'sports', 'self defense'], category: 'Sports & Fitness' },
	{ emoji: '🤺', keywords: ['fencing', 'sports', 'athletics', 'swordsmanship', 'olympic'], category: 'Sports & Fitness' },
	{ emoji: '♟️', keywords: ['chess', 'strategy', 'games', 'logic', 'competition'], category: 'Sports & Fitness' },
	{ emoji: '🎲', keywords: ['games', 'probability', 'statistics', 'board games', 'chance'], category: 'Sports & Fitness' },
	{ emoji: '🎮', keywords: ['gaming', 'video games', 'game design', 'esports', 'interactive'], category: 'Sports & Fitness' },
	{ emoji: '🏆', keywords: ['competition', 'awards', 'sports', 'achievement', 'championship'], category: 'Sports & Fitness' },
	{ emoji: '🥇', keywords: ['competition', 'achievement', 'sports', 'first place', 'excellence'], category: 'Sports & Fitness' },

	// Education & Academic
	{ emoji: '🎓', keywords: ['graduation', 'academic', 'degree', 'university', 'college', 'exam'], category: 'Education & Academic' },
	{ emoji: '🏫', keywords: ['school', 'education', 'learning', 'teaching', 'academic'], category: 'Education & Academic' },
	{ emoji: '👩‍🏫', keywords: ['teaching', 'education', 'instructor', 'professor', 'lecture'], category: 'Education & Academic' },
	{ emoji: '👨‍🎓', keywords: ['student', 'education', 'university', 'college', 'learning'], category: 'Education & Academic' },
	{ emoji: '🎒', keywords: ['school', 'student', 'education', 'backpack', 'supplies'], category: 'Education & Academic' },
	{ emoji: '📏', keywords: ['measurement', 'math', 'geometry', 'ruler', 'precision'], category: 'Education & Academic' },
	{ emoji: '🔎', keywords: ['research', 'investigation', 'analysis', 'study', 'search'], category: 'Education & Academic' },
	{ emoji: '🔍', keywords: ['research', 'investigation', 'analysis', 'study', 'detail'], category: 'Education & Academic' },

	// Travel & Geography
	{ emoji: '🗺️', keywords: ['geography', 'maps', 'travel', 'cartography', 'exploration'], category: 'Travel & Geography' },
	{ emoji: '🧭', keywords: ['navigation', 'geography', 'compass', 'orientation', 'exploration'], category: 'Travel & Geography' },
	{ emoji: '🗽', keywords: ['american studies', 'history', 'travel', 'landmarks', 'usa'], category: 'Travel & Geography' },
	{ emoji: '🗼', keywords: ['architecture', 'landmarks', 'travel', 'paris', 'structures'], category: 'Travel & Geography' },
	{ emoji: '🏰', keywords: ['history', 'castles', 'medieval', 'architecture', 'european'], category: 'Travel & Geography' },
	{ emoji: '⛩️', keywords: ['japanese', 'culture', 'religion', 'architecture', 'asian studies'], category: 'Travel & Geography' },
	{ emoji: '🕌', keywords: ['islamic studies', 'architecture', 'religion', 'culture', 'mosque'], category: 'Travel & Geography' },
	{ emoji: '⛪', keywords: ['religion', 'christianity', 'architecture', 'church', 'theology'], category: 'Travel & Geography' },
	{ emoji: '🛕', keywords: ['hinduism', 'temple', 'religion', 'architecture', 'indian'], category: 'Travel & Geography' },

	// Mythology & Fantasy
	{ emoji: '🐉', keywords: ['mythology', 'dragons', 'fantasy', 'legends', 'folklore'], category: 'Mythology & Fantasy' },
	{ emoji: '🦄', keywords: ['mythology', 'unicorn', 'fantasy', 'legends', 'fairy tales'], category: 'Mythology & Fantasy' },
	{ emoji: '🧙', keywords: ['fantasy', 'magic', 'wizard', 'mythology', 'folklore'], category: 'Mythology & Fantasy' },
	{ emoji: '🧚', keywords: ['fantasy', 'fairy', 'mythology', 'folklore', 'fairy tales'], category: 'Mythology & Fantasy' },
	{ emoji: '🧜', keywords: ['mythology', 'mermaids', 'fantasy', 'folklore', 'ocean'], category: 'Mythology & Fantasy' },
	{ emoji: '🪄', keywords: ['magic', 'fantasy', 'illusion', 'tricks', 'wonder'], category: 'Mythology & Fantasy' },
	{ emoji: '✨', keywords: ['magic', 'sparkle', 'special', 'creativity', 'highlights'], category: 'Mythology & Fantasy' },

	// Military & Defense
	{ emoji: '🎖️', keywords: ['military', 'medals', 'defense', 'honors', 'service'], category: 'Military & Defense' },
	{ emoji: '🛡️', keywords: ['defense', 'security', 'protection', 'military', 'shields'], category: 'Military & Defense' },
	{ emoji: '⚔️', keywords: ['military history', 'combat', 'swords', 'warfare', 'medieval'], category: 'Military & Defense' },
	{ emoji: '🪖', keywords: ['military', 'army', 'defense', 'helmet', 'soldier'], category: 'Military & Defense' },

	// Miscellaneous
	{ emoji: '🔥', keywords: ['trending', 'hot topics', 'popular', 'fire', 'urgent'], category: 'Miscellaneous' },
	{ emoji: '💫', keywords: ['highlights', 'special', 'important', 'dizzy', 'stars'], category: 'Miscellaneous' },
	{ emoji: '🌈', keywords: ['diversity', 'lgbtq', 'pride', 'colors', 'inclusivity'], category: 'Miscellaneous' },
	{ emoji: '🎪', keywords: ['circus', 'entertainment', 'performing arts', 'events', 'shows'], category: 'Miscellaneous' },
	{ emoji: '🎠', keywords: ['amusement', 'recreation', 'entertainment', 'fun', 'parks'], category: 'Miscellaneous' },
	{ emoji: '🎡', keywords: ['amusement', 'recreation', 'entertainment', 'ferris wheel', 'parks'], category: 'Miscellaneous' },
	{ emoji: '🎢', keywords: ['amusement', 'recreation', 'entertainment', 'roller coaster', 'physics'], category: 'Miscellaneous' },
	{ emoji: '🌺', keywords: ['hawaiian', 'tropical', 'flowers', 'culture', 'nature'], category: 'Miscellaneous' },
	{ emoji: '🐕‍🦺', keywords: ['service animals', 'assistance', 'disability studies', 'therapy', 'support'], category: 'Miscellaneous' },
	{ emoji: '🤟', keywords: ['sign language', 'deaf studies', 'communication', 'asl', 'accessibility'], category: 'Miscellaneous' },
	{ emoji: '♿', keywords: ['disability studies', 'accessibility', 'inclusion', 'accommodation', 'ada'], category: 'Miscellaneous' },
];

// Get unique categories in display order
export const emojiCategories = [
	'Books & Notebooks',
	'Sciences',
	'Math & Logic',
	'Technology & Engineering',
	'Medicine & Health',
	'Humanities & Social Sciences',
	'Religion & Philosophy',
	'Business & Economics',
	'Arts & Design',
	'Animals & Zoology',
	'Agriculture & Environment',
	'Food & Culinary',
	'Sports & Fitness',
	'Education & Academic',
	'Travel & Geography',
	'Mythology & Fantasy',
	'Military & Defense',
	'Miscellaneous',
] as const;

export type EmojiCategory = (typeof emojiCategories)[number];

// Group emojis by category for display
export function getEmojisByCategory(): Map<EmojiCategory, NotebookEmoji[]> {
	const map = new Map<EmojiCategory, NotebookEmoji[]>();
	for (const category of emojiCategories) {
		map.set(category, notebookEmojis.filter((e) => e.category === category));
	}
	return map;
}

// Search emojis by keyword
export function searchEmojis(query: string): NotebookEmoji[] {
	const lowerQuery = query.toLowerCase().trim();
	if (!lowerQuery) return notebookEmojis;
	return notebookEmojis.filter(
		(e) =>
			e.emoji.includes(query) ||
			e.keywords.some((k) => k.toLowerCase().includes(lowerQuery))
	);
}
