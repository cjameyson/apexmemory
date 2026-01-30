import React, { useState, useEffect, useRef } from 'react';
import { Book, FileText, ChevronDown, Plus, Search, Mic, Youtube, Link, StickyNote, Command, Star, ArrowRight, Check, Zap, X, Send, Upload, MessageSquare, ChevronRight, Eye, Box, Play, BookOpen, Grid, ChevronLeft, Tag, MoreHorizontal, Volume2, ExternalLink, Info, Edit2, Home, TrendingUp, Calendar, Target, Award, Clock, BarChart3, Layers, Mouse, Type, Image, Scissors, Edit, ArrowUpRight, Maximize2, Minimize2, SkipBack, SkipForward, Pause, RotateCcw, ZoomIn, ZoomOut, Download, Bookmark, List, RefreshCw, Menu, ChevronsLeft, ChevronsRight, Filter, Trash2, Settings, FolderOpen, Copy, Archive, Table, LayoutGrid, Hash, AlertCircle, Square, CheckSquare } from 'lucide-react';

// ============================================================================
// MOCK DATA
// ============================================================================

const mockNotebooks = [
  { id: 1, name: 'Biology 101', emoji: '🧬', color: '#0ea5e9', dueCount: 23, streak: 7, totalCards: 156, totalFacts: 89, retention: 87 },
  { id: 2, name: 'Spanish B2', emoji: '🇪🇸', color: '#f59e0b', dueCount: 45, streak: 12, totalCards: 412, totalFacts: 156, retention: 82 },
  { id: 3, name: 'Machine Learning', emoji: '🤖', color: '#6366f1', dueCount: 12, streak: 3, totalCards: 89, totalFacts: 42, retention: 91 },
  { id: 4, name: 'History - WWII', emoji: '📜', color: '#ec4899', dueCount: 8, streak: 5, totalCards: 234, totalFacts: 98, retention: 79 },
];

const mockSources = {
  1: [
    { id: 's1', name: 'Chapter 3 - Cell Division', type: 'pdf', cards: 24, excerpt: 'Cell division is the process...', pages: 42, addedAt: '2 days ago' },
    { id: 's2', name: 'Mitosis vs Meiosis', type: 'youtube', cards: 18, excerpt: 'A comprehensive comparison...', duration: '12:34', addedAt: '1 week ago' },
    { id: 's3', name: 'Lab Notes - Microscopy', type: 'notes', cards: 8, excerpt: 'Observations from the microscopy lab...', pages: 5, addedAt: '3 days ago' },
  ],
};

// Facts mock data - representing the new schema structure
const mockFacts = {
  1: [
    {
      id: 'f1',
      factType: 'basic',
      content: {
        version: 1,
        fields: [
          { name: 'front', type: 'plain_text', value: 'What is the powerhouse of the cell?' },
          { name: 'back', type: 'rich_text', value: 'The mitochondria - responsible for producing ATP through cellular respiration' }
        ]
      },
      sourceId: 's1',
      cardCount: 1,
      dueCount: 1,
      tags: ['organelles', 'exam-1'],
      createdAt: '2024-01-15T10:30:00Z',
      updatedAt: '2024-01-20T14:22:00Z'
    },
    {
      id: 'f2',
      factType: 'basic',
      content: {
        version: 1,
        fields: [
          { name: 'front', type: 'plain_text', value: 'What are the phases of mitosis?' },
          { name: 'back', type: 'rich_text', value: 'Prophase, Metaphase, Anaphase, Telophase (PMAT)' }
        ]
      },
      sourceId: 's1',
      cardCount: 1,
      dueCount: 1,
      tags: ['cell-division', 'exam-1'],
      createdAt: '2024-01-15T10:35:00Z',
      updatedAt: '2024-01-18T09:15:00Z'
    },
    {
      id: 'f3',
      factType: 'cloze',
      content: {
        version: 1,
        fields: [
          { name: 'text', type: 'cloze_text', value: '{{c1::Mitochondria}} are the {{c2::powerhouse}} of the cell, producing {{c3::ATP}} through cellular respiration.' },
          { name: 'extra', type: 'rich_text', value: 'Remember: ATP = Adenosine Triphosphate' }
        ]
      },
      sourceId: 's1',
      cardCount: 3,
      dueCount: 2,
      tags: ['organelles', 'energy'],
      createdAt: '2024-01-16T11:00:00Z',
      updatedAt: '2024-01-21T16:45:00Z'
    },
    {
      id: 'f4',
      factType: 'cloze',
      content: {
        version: 1,
        fields: [
          { name: 'text', type: 'cloze_text', value: 'The cell cycle consists of {{c1::Interphase}} and {{c2::M phase (Mitosis)}}. Interphase includes {{c3::G1}}, {{c4::S}}, and {{c5::G2}} phases.' },
          { name: 'extra', type: 'rich_text', value: 'G1: Growth, S: DNA Synthesis, G2: Preparation for division' }
        ]
      },
      sourceId: 's2',
      cardCount: 5,
      dueCount: 3,
      tags: ['cell-cycle', 'exam-1'],
      createdAt: '2024-01-17T09:20:00Z',
      updatedAt: '2024-01-22T11:30:00Z'
    },
    {
      id: 'f5',
      factType: 'image_occlusion',
      content: {
        version: 1,
        fields: [
          { name: 'title', type: 'plain_text', value: 'Animal Cell Diagram' },
          { name: 'image', type: 'image', value: '/images/animal-cell.png' },
          { name: 'masks', type: 'masks', value: [
            { id: 'm_abc123', label: 'Nucleus', x: 150, y: 120, w: 80, h: 80 },
            { id: 'm_def456', label: 'Mitochondria', x: 280, y: 200, w: 60, h: 40 },
            { id: 'm_ghi789', label: 'Endoplasmic Reticulum', x: 100, y: 200, w: 90, h: 50 },
            { id: 'm_jkl012', label: 'Golgi Apparatus', x: 320, y: 150, w: 70, h: 45 }
          ]}
        ]
      },
      sourceId: 's3',
      cardCount: 4,
      dueCount: 1,
      tags: ['cell-structure', 'diagrams'],
      createdAt: '2024-01-18T14:00:00Z',
      updatedAt: '2024-01-23T10:00:00Z'
    },
    {
      id: 'f6',
      factType: 'basic',
      content: {
        version: 1,
        fields: [
          { name: 'front', type: 'plain_text', value: 'Define apoptosis' },
          { name: 'back', type: 'rich_text', value: 'Programmed cell death - a controlled process of eliminating unwanted or damaged cells' }
        ]
      },
      sourceId: 's2',
      cardCount: 1,
      dueCount: 0,
      tags: ['cell-death'],
      createdAt: '2024-01-14T08:00:00Z',
      updatedAt: '2024-01-19T12:00:00Z'
    },
    {
      id: 'f7',
      factType: 'cloze',
      content: {
        version: 1,
        fields: [
          { name: 'text', type: 'cloze_text', value: 'DNA replication is {{c1::semi-conservative}}, meaning each new double helix contains {{c2::one original strand}} and {{c3::one newly synthesized strand}}.' },
          { name: 'extra', type: 'rich_text', value: '' }
        ]
      },
      sourceId: null,
      cardCount: 3,
      dueCount: 2,
      tags: ['dna', 'replication'],
      createdAt: '2024-01-19T16:30:00Z',
      updatedAt: '2024-01-24T09:00:00Z'
    },
    {
      id: 'f8',
      factType: 'image_occlusion',
      content: {
        version: 1,
        fields: [
          { name: 'title', type: 'plain_text', value: 'Phases of Mitosis' },
          { name: 'image', type: 'image', value: '/images/mitosis-phases.png' },
          { name: 'masks', type: 'masks', value: [
            { id: 'm_mit001', label: 'Prophase', x: 50, y: 100, w: 100, h: 80 },
            { id: 'm_mit002', label: 'Metaphase', x: 180, y: 100, w: 100, h: 80 },
            { id: 'm_mit003', label: 'Anaphase', x: 310, y: 100, w: 100, h: 80 },
            { id: 'm_mit004', label: 'Telophase', x: 440, y: 100, w: 100, h: 80 }
          ]}
        ]
      },
      sourceId: 's1',
      cardCount: 4,
      dueCount: 4,
      tags: ['mitosis', 'exam-1', 'diagrams'],
      createdAt: '2024-01-20T10:00:00Z',
      updatedAt: '2024-01-25T08:30:00Z'
    },
    {
      id: 'f9',
      factType: 'basic',
      content: {
        version: 1,
        fields: [
          { name: 'front', type: 'plain_text', value: 'What is crossing over?' },
          { name: 'back', type: 'rich_text', value: 'Exchange of genetic material between homologous chromosomes during meiosis I, increasing genetic diversity' }
        ]
      },
      sourceId: 's1',
      cardCount: 1,
      dueCount: 1,
      tags: ['meiosis', 'genetics'],
      createdAt: '2024-01-21T11:30:00Z',
      updatedAt: '2024-01-26T14:00:00Z'
    },
    {
      id: 'f10',
      factType: 'basic',
      content: {
        version: 1,
        fields: [
          { name: 'front', type: 'plain_text', value: 'How many chromosomes do human somatic cells have?' },
          { name: 'back', type: 'rich_text', value: '46 (23 pairs) - called the diploid number (2n)' }
        ]
      },
      sourceId: 's3',
      cardCount: 1,
      dueCount: 0,
      tags: ['genetics', 'chromosomes'],
      createdAt: '2024-01-22T09:00:00Z',
      updatedAt: '2024-01-27T10:30:00Z'
    },
    {
      id: 'f11',
      factType: 'cloze',
      content: {
        version: 1,
        fields: [
          { name: 'text', type: 'cloze_text', value: 'The enzyme {{c1::helicase}} unwinds DNA by breaking {{c2::hydrogen bonds}} between base pairs.' },
          { name: 'extra', type: 'rich_text', value: '' }
        ]
      },
      sourceId: null,
      cardCount: 2,
      dueCount: 1,
      tags: ['dna', 'enzymes'],
      createdAt: '2024-01-23T13:45:00Z',
      updatedAt: '2024-01-28T11:00:00Z'
    },
    {
      id: 'f12',
      factType: 'basic',
      content: {
        version: 1,
        fields: [
          { name: 'front', type: 'plain_text', value: 'Difference between chromatin and chromosomes?' },
          { name: 'back', type: 'rich_text', value: 'Chromatin is loosely coiled DNA (during interphase); chromosomes are tightly condensed chromatin (during division)' }
        ]
      },
      sourceId: 's1',
      cardCount: 1,
      dueCount: 0,
      tags: ['cell-division', 'dna'],
      createdAt: '2024-01-24T15:00:00Z',
      updatedAt: '2024-01-29T09:30:00Z'
    },
  ],
};

// Cards for each fact
const mockCards = {
  'f1': [
    { id: 'c1', elementId: '', state: 'review', due: '2024-01-30T10:00:00Z', interval: '1d', stability: 2.5, difficulty: 5.2, reps: 4, lapses: 1 }
  ],
  'f2': [
    { id: 'c2', elementId: '', state: 'learning', due: '2024-01-29T14:00:00Z', interval: '4h', stability: null, difficulty: null, reps: 2, lapses: 0, step: 1 }
  ],
  'f3': [
    { id: 'c3a', elementId: 'c1', state: 'review', due: '2024-01-28T09:00:00Z', interval: '3d', stability: 4.1, difficulty: 4.8, reps: 6, lapses: 0 },
    { id: 'c3b', elementId: 'c2', state: 'review', due: '2024-01-30T16:00:00Z', interval: '2d', stability: 3.2, difficulty: 5.5, reps: 5, lapses: 1 },
    { id: 'c3c', elementId: 'c3', state: 'new', due: null, interval: null, stability: null, difficulty: null, reps: 0, lapses: 0 }
  ],
  'f4': [
    { id: 'c4a', elementId: 'c1', state: 'review', due: '2024-01-29T11:00:00Z', interval: '1d', stability: 2.1, difficulty: 6.0, reps: 3, lapses: 2 },
    { id: 'c4b', elementId: 'c2', state: 'review', due: '2024-01-30T08:00:00Z', interval: '1d', stability: 2.3, difficulty: 5.8, reps: 3, lapses: 1 },
    { id: 'c4c', elementId: 'c3', state: 'learning', due: '2024-01-29T15:30:00Z', interval: '10m', stability: null, difficulty: null, reps: 1, lapses: 0, step: 0 },
    { id: 'c4d', elementId: 'c4', state: 'new', due: null, interval: null, stability: null, difficulty: null, reps: 0, lapses: 0 },
    { id: 'c4e', elementId: 'c5', state: 'new', due: null, interval: null, stability: null, difficulty: null, reps: 0, lapses: 0 }
  ],
  'f5': [
    { id: 'c5a', elementId: 'm_abc123', state: 'review', due: '2024-01-31T10:00:00Z', interval: '4d', stability: 5.2, difficulty: 4.2, reps: 8, lapses: 0 },
    { id: 'c5b', elementId: 'm_def456', state: 'review', due: '2024-02-02T10:00:00Z', interval: '6d', stability: 6.8, difficulty: 3.9, reps: 10, lapses: 0 },
    { id: 'c5c', elementId: 'm_ghi789', state: 'review', due: '2024-02-01T10:00:00Z', interval: '5d', stability: 5.9, difficulty: 4.5, reps: 9, lapses: 1 },
    { id: 'c5d', elementId: 'm_jkl012', state: 'relearning', due: '2024-01-29T16:00:00Z', interval: '1h', stability: 1.2, difficulty: 6.5, reps: 7, lapses: 2, step: 0 }
  ],
};

const mockGlobalStats = {
  totalCards: 891,
  cardsReviewedToday: 34,
  currentStreak: 12,
  longestStreak: 28,
  averageRetention: 85,
  reviewsThisWeek: [12, 45, 38, 52, 34, 0, 0],
  upcomingDue: { today: 88, tomorrow: 42, thisWeek: 156 },
};

// ============================================================================
// UTILITY COMPONENTS
// ============================================================================

const SourceIcon = ({ type, className }) => {
  const icons = { pdf: <FileText className={className} />, youtube: <Youtube className={className} />, url: <Link className={className} />, audio: <Mic className={className} />, notes: <StickyNote className={className} /> };
  return icons[type] || <FileText className={className} />;
};

const FactTypeIcon = ({ type, className }) => {
  const icons = {
    basic: <Layers className={className} />,
    cloze: <Type className={className} />,
    image_occlusion: <Image className={className} />
  };
  return icons[type] || <Layers className={className} />;
};

const FactTypeBadge = ({ type }) => {
  const styles = {
    basic: 'bg-slate-100 text-slate-600',
    cloze: 'bg-purple-100 text-purple-600',
    image_occlusion: 'bg-amber-100 text-amber-600'
  };
  const labels = {
    basic: 'Basic',
    cloze: 'Cloze',
    image_occlusion: 'Image'
  };
  return (
    <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium ${styles[type] || styles.basic}`}>
      <FactTypeIcon type={type} className="w-3 h-3" />
      {labels[type] || type}
    </span>
  );
};

const CardStateBadge = ({ state, step }) => {
  const styles = {
    new: 'bg-blue-100 text-blue-600',
    learning: 'bg-amber-100 text-amber-600',
    review: 'bg-emerald-100 text-emerald-600',
    relearning: 'bg-red-100 text-red-600'
  };
  return (
    <span className={`px-1.5 py-0.5 rounded text-xs font-medium ${styles[state]}`}>
      {state}{step !== undefined && step !== null ? ` (${step})` : ''}
    </span>
  );
};

const ProgressRing = ({ progress, size = 40, stroke = 4, trackColor = '#e2e8f0', progressColor = '#0ea5e9' }) => {
  const radius = (size - stroke) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (progress / 100) * circumference;
  return (
    <svg width={size} height={size} className="transform -rotate-90">
      <circle cx={size / 2} cy={size / 2} r={radius} fill="none" stroke={trackColor} strokeWidth={stroke} />
      <circle cx={size / 2} cy={size / 2} r={radius} fill="none" stroke={progressColor} strokeWidth={stroke} strokeDasharray={circumference} strokeDashoffset={offset} strokeLinecap="round" className="transition-all duration-500" />
    </svg>
  );
};

// ============================================================================
// FACT DISPLAY HELPERS
// ============================================================================

// Get the primary display text for a fact based on its type
const getFactDisplayText = (fact) => {
  const { factType, content } = fact;
  const fields = content.fields;
  
  switch (factType) {
    case 'basic': {
      const frontField = fields.find(f => f.name === 'front');
      return frontField?.value || '';
    }
    case 'cloze': {
      const textField = fields.find(f => f.name === 'text');
      if (!textField?.value) return '';
      // Replace cloze deletions with blanks: {{c1::text}} → [...]
      return textField.value.replace(/\{\{c\d+::([^}]+)\}\}/g, '[...]');
    }
    case 'image_occlusion': {
      const titleField = fields.find(f => f.name === 'title');
      return titleField?.value || 'Untitled Image';
    }
    default:
      return '';
  }
};

// Get the answer/back text for a fact (for preview)
const getFactAnswerText = (fact) => {
  const { factType, content } = fact;
  const fields = content.fields;
  
  switch (factType) {
    case 'basic': {
      const backField = fields.find(f => f.name === 'back');
      return backField?.value || '';
    }
    case 'cloze': {
      const textField = fields.find(f => f.name === 'text');
      if (!textField?.value) return '';
      // Show full text with cloze content visible
      return textField.value.replace(/\{\{c\d+::([^}]+)\}\}/g, '$1');
    }
    case 'image_occlusion': {
      const masksField = fields.find(f => f.name === 'masks');
      const masks = masksField?.value || [];
      return masks.map(m => m.label).join(', ');
    }
    default:
      return '';
  }
};

// ============================================================================
// FACT BROWSER COMPONENT
// ============================================================================

const FactBrowser = ({ 
  facts, 
  sources, 
  notebook,
  onCreateFact,
  onEditFact,
  onDeleteFact,
  onStartReview
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [typeFilter, setTypeFilter] = useState('all');
  const [sortBy, setSortBy] = useState('updated');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [expandedFact, setExpandedFact] = useState(null);
  const [selectedFacts, setSelectedFacts] = useState(new Set());
  const [viewMode, setViewMode] = useState('table'); // 'table' or 'grid'
  
  const PAGE_SIZE = 10;
  
  // Filter and sort facts
  const filteredFacts = facts.filter(fact => {
    // Type filter
    if (typeFilter !== 'all' && fact.factType !== typeFilter) return false;
    
    // Search filter
    if (searchQuery) {
      const displayText = getFactDisplayText(fact).toLowerCase();
      const answerText = getFactAnswerText(fact).toLowerCase();
      const query = searchQuery.toLowerCase();
      if (!displayText.includes(query) && !answerText.includes(query)) return false;
    }
    
    return true;
  });
  
  // Sort
  const sortedFacts = [...filteredFacts].sort((a, b) => {
    let comparison = 0;
    switch (sortBy) {
      case 'updated':
        comparison = new Date(b.updatedAt) - new Date(a.updatedAt);
        break;
      case 'created':
        comparison = new Date(b.createdAt) - new Date(a.createdAt);
        break;
      case 'cards':
        comparison = b.cardCount - a.cardCount;
        break;
      case 'due':
        comparison = b.dueCount - a.dueCount;
        break;
      default:
        comparison = 0;
    }
    return sortOrder === 'desc' ? comparison : -comparison;
  });
  
  // Pagination
  const totalPages = Math.ceil(sortedFacts.length / PAGE_SIZE);
  const paginatedFacts = sortedFacts.slice((currentPage - 1) * PAGE_SIZE, currentPage * PAGE_SIZE);
  
  // Stats
  const totalCards = facts.reduce((sum, f) => sum + f.cardCount, 0);
  const totalDue = facts.reduce((sum, f) => sum + f.dueCount, 0);
  const basicCount = facts.filter(f => f.factType === 'basic').length;
  const clozeCount = facts.filter(f => f.factType === 'cloze').length;
  const imageCount = facts.filter(f => f.factType === 'image_occlusion').length;
  
  // Selection handlers
  const toggleSelectAll = () => {
    if (selectedFacts.size === paginatedFacts.length) {
      setSelectedFacts(new Set());
    } else {
      setSelectedFacts(new Set(paginatedFacts.map(f => f.id)));
    }
  };
  
  const toggleSelect = (factId) => {
    const newSelected = new Set(selectedFacts);
    if (newSelected.has(factId)) {
      newSelected.delete(factId);
    } else {
      newSelected.add(factId);
    }
    setSelectedFacts(newSelected);
  };
  
  return (
    <div className="flex-1 flex flex-col bg-white overflow-hidden">
      {/* Header with Stats */}
      <div className="px-6 py-4 border-b border-slate-200 bg-gradient-to-r from-slate-50 to-white">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-xl font-bold text-slate-900">Facts & Cards</h2>
            <p className="text-sm text-slate-500 mt-0.5">
              {facts.length} facts • {totalCards} cards • {totalDue} due for review
            </p>
          </div>
          <div className="flex items-center gap-2">
            {totalDue > 0 && (
              <button 
                onClick={() => onStartReview && onStartReview()}
                className="flex items-center gap-2 px-4 py-2 bg-sky-500 hover:bg-sky-600 rounded-lg text-white text-sm font-medium transition-colors"
              >
                <Play className="w-4 h-4" />
                Review ({totalDue})
              </button>
            )}
            <button 
              onClick={() => onCreateFact && onCreateFact()}
              className="flex items-center gap-2 px-4 py-2 bg-slate-900 hover:bg-slate-800 rounded-lg text-white text-sm font-medium transition-colors"
            >
              <Plus className="w-4 h-4" />
              Create Fact
            </button>
          </div>
        </div>
        
        {/* Quick Stats */}
        <div className="flex items-center gap-6">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-slate-100 rounded-lg flex items-center justify-center">
              <Layers className="w-4 h-4 text-slate-500" />
            </div>
            <div>
              <p className="text-xs text-slate-500">Basic</p>
              <p className="text-sm font-semibold text-slate-900">{basicCount}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-purple-100 rounded-lg flex items-center justify-center">
              <Type className="w-4 h-4 text-purple-500" />
            </div>
            <div>
              <p className="text-xs text-slate-500">Cloze</p>
              <p className="text-sm font-semibold text-slate-900">{clozeCount}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-amber-100 rounded-lg flex items-center justify-center">
              <Image className="w-4 h-4 text-amber-500" />
            </div>
            <div>
              <p className="text-xs text-slate-500">Image</p>
              <p className="text-sm font-semibold text-slate-900">{imageCount}</p>
            </div>
          </div>
          <div className="w-px h-8 bg-slate-200" />
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-sky-100 rounded-lg flex items-center justify-center">
              <Clock className="w-4 h-4 text-sky-500" />
            </div>
            <div>
              <p className="text-xs text-slate-500">Due Today</p>
              <p className="text-sm font-semibold text-sky-600">{totalDue}</p>
            </div>
          </div>
        </div>
      </div>
      
      {/* Toolbar */}
      <div className="px-6 py-3 border-b border-slate-200 flex items-center justify-between gap-4">
        <div className="flex items-center gap-3 flex-1">
          {/* Search */}
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              type="text"
              placeholder="Search facts..."
              value={searchQuery}
              onChange={(e) => { setSearchQuery(e.target.value); setCurrentPage(1); }}
              className="w-full pl-10 pr-4 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent"
            />
          </div>
          
          {/* Type Filter */}
          <div className="flex items-center bg-slate-100 rounded-lg p-0.5">
            {[
              { id: 'all', label: 'All' },
              { id: 'basic', label: 'Basic', icon: Layers },
              { id: 'cloze', label: 'Cloze', icon: Type },
              { id: 'image_occlusion', label: 'Image', icon: Image },
            ].map(filter => (
              <button
                key={filter.id}
                onClick={() => { setTypeFilter(filter.id); setCurrentPage(1); }}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium transition-all ${
                  typeFilter === filter.id 
                    ? 'bg-white shadow-sm text-slate-900' 
                    : 'text-slate-500 hover:text-slate-700'
                }`}
              >
                {filter.icon && <filter.icon className="w-3.5 h-3.5" />}
                {filter.label}
              </button>
            ))}
          </div>
        </div>
        
        <div className="flex items-center gap-2">
          {/* Sort */}
          <select
            value={`${sortBy}-${sortOrder}`}
            onChange={(e) => {
              const [field, order] = e.target.value.split('-');
              setSortBy(field);
              setSortOrder(order);
            }}
            className="px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-sky-500"
          >
            <option value="updated-desc">Recently Updated</option>
            <option value="updated-asc">Oldest Updated</option>
            <option value="created-desc">Recently Created</option>
            <option value="created-asc">Oldest Created</option>
            <option value="cards-desc">Most Cards</option>
            <option value="due-desc">Most Due</option>
          </select>
          
          {/* View Mode Toggle */}
          <div className="flex items-center bg-slate-100 rounded-lg p-0.5">
            <button
              onClick={() => setViewMode('table')}
              className={`p-2 rounded-md transition-all ${viewMode === 'table' ? 'bg-white shadow-sm' : 'hover:bg-slate-200'}`}
              title="Table view"
            >
              <Table className={`w-4 h-4 ${viewMode === 'table' ? 'text-slate-900' : 'text-slate-500'}`} />
            </button>
            <button
              onClick={() => setViewMode('grid')}
              className={`p-2 rounded-md transition-all ${viewMode === 'grid' ? 'bg-white shadow-sm' : 'hover:bg-slate-200'}`}
              title="Grid view"
            >
              <LayoutGrid className={`w-4 h-4 ${viewMode === 'grid' ? 'text-slate-900' : 'text-slate-500'}`} />
            </button>
          </div>
        </div>
      </div>
      
      {/* Bulk Actions Bar (when items selected) */}
      {selectedFacts.size > 0 && (
        <div className="px-6 py-2 bg-sky-50 border-b border-sky-100 flex items-center justify-between">
          <div className="flex items-center gap-2 text-sm">
            <span className="font-medium text-sky-700">{selectedFacts.size} selected</span>
            <button 
              onClick={() => setSelectedFacts(new Set())}
              className="text-sky-600 hover:text-sky-800 underline"
            >
              Clear
            </button>
          </div>
          <div className="flex items-center gap-2">
            <button className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-slate-600 hover:bg-sky-100 rounded-lg transition-colors">
              <Tag className="w-4 h-4" />
              Add Tags
            </button>
            <button className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-slate-600 hover:bg-sky-100 rounded-lg transition-colors">
              <Archive className="w-4 h-4" />
              Suspend
            </button>
            <button className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-red-600 hover:bg-red-50 rounded-lg transition-colors">
              <Trash2 className="w-4 h-4" />
              Delete
            </button>
          </div>
        </div>
      )}
      
      {/* Content Area */}
      <div className="flex-1 overflow-y-auto">
        {viewMode === 'table' ? (
          <FactTable
            facts={paginatedFacts}
            sources={sources}
            selectedFacts={selectedFacts}
            expandedFact={expandedFact}
            onToggleSelect={toggleSelect}
            onToggleSelectAll={toggleSelectAll}
            onToggleExpand={(id) => setExpandedFact(expandedFact === id ? null : id)}
            onEdit={onEditFact}
            onDelete={onDeleteFact}
            allSelected={selectedFacts.size === paginatedFacts.length && paginatedFacts.length > 0}
          />
        ) : (
          <FactGrid
            facts={paginatedFacts}
            sources={sources}
            selectedFacts={selectedFacts}
            onToggleSelect={toggleSelect}
            onEdit={onEditFact}
            onDelete={onDeleteFact}
          />
        )}
        
        {/* Empty State */}
        {filteredFacts.length === 0 && (
          <div className="flex flex-col items-center justify-center py-16 px-4">
            <div className="w-16 h-16 bg-slate-100 rounded-2xl flex items-center justify-center mb-4">
              <Layers className="w-8 h-8 text-slate-400" />
            </div>
            <h3 className="text-lg font-semibold text-slate-900 mb-1">
              {searchQuery ? 'No facts found' : 'No facts yet'}
            </h3>
            <p className="text-slate-500 text-sm text-center max-w-sm mb-6">
              {searchQuery 
                ? `No facts match "${searchQuery}". Try a different search term.`
                : 'Create your first fact to start building your knowledge base.'
              }
            </p>
            {!searchQuery && (
              <button 
                onClick={() => onCreateFact && onCreateFact()}
                className="flex items-center gap-2 px-4 py-2 bg-sky-500 hover:bg-sky-600 rounded-lg text-white text-sm font-medium transition-colors"
              >
                <Plus className="w-4 h-4" />
                Create your first fact
              </button>
            )}
          </div>
        )}
      </div>
      
      {/* Pagination */}
      {totalPages > 1 && (
        <div className="px-6 py-3 border-t border-slate-200 flex items-center justify-between">
          <p className="text-sm text-slate-500">
            Showing {(currentPage - 1) * PAGE_SIZE + 1} to {Math.min(currentPage * PAGE_SIZE, filteredFacts.length)} of {filteredFacts.length} facts
          </p>
          <div className="flex items-center gap-1">
            <button
              onClick={() => setCurrentPage(1)}
              disabled={currentPage === 1}
              className="p-2 rounded-lg hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronsLeft className="w-4 h-4 text-slate-600" />
            </button>
            <button
              onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
              disabled={currentPage === 1}
              className="p-2 rounded-lg hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronLeft className="w-4 h-4 text-slate-600" />
            </button>
            
            {/* Page numbers */}
            <div className="flex items-center gap-1 mx-2">
              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                let pageNum;
                if (totalPages <= 5) {
                  pageNum = i + 1;
                } else if (currentPage <= 3) {
                  pageNum = i + 1;
                } else if (currentPage >= totalPages - 2) {
                  pageNum = totalPages - 4 + i;
                } else {
                  pageNum = currentPage - 2 + i;
                }
                return (
                  <button
                    key={pageNum}
                    onClick={() => setCurrentPage(pageNum)}
                    className={`w-8 h-8 rounded-lg text-sm font-medium transition-colors ${
                      currentPage === pageNum
                        ? 'bg-sky-500 text-white'
                        : 'hover:bg-slate-100 text-slate-600'
                    }`}
                  >
                    {pageNum}
                  </button>
                );
              })}
            </div>
            
            <button
              onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
              disabled={currentPage === totalPages}
              className="p-2 rounded-lg hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronRight className="w-4 h-4 text-slate-600" />
            </button>
            <button
              onClick={() => setCurrentPage(totalPages)}
              disabled={currentPage === totalPages}
              className="p-2 rounded-lg hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronsRight className="w-4 h-4 text-slate-600" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

// ============================================================================
// FACT TABLE COMPONENT
// ============================================================================

const FactTable = ({ 
  facts, 
  sources, 
  selectedFacts, 
  expandedFact,
  onToggleSelect, 
  onToggleSelectAll,
  onToggleExpand,
  onEdit, 
  onDelete,
  allSelected
}) => {
  return (
    <table className="w-full">
      <thead className="sticky top-0 bg-slate-50 z-10">
        <tr className="border-b border-slate-200">
          <th className="w-12 px-4 py-3 text-left">
            <button onClick={onToggleSelectAll} className="p-1 hover:bg-slate-200 rounded transition-colors">
              {allSelected ? (
                <CheckSquare className="w-4 h-4 text-sky-500" />
              ) : (
                <Square className="w-4 h-4 text-slate-400" />
              )}
            </button>
          </th>
          <th className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">
            Content
          </th>
          <th className="w-24 px-4 py-3 text-center text-xs font-medium text-slate-500 uppercase tracking-wider">
            Type
          </th>
          <th className="w-20 px-4 py-3 text-center text-xs font-medium text-slate-500 uppercase tracking-wider">
            Cards
          </th>
          <th className="w-20 px-4 py-3 text-center text-xs font-medium text-slate-500 uppercase tracking-wider">
            Due
          </th>
          <th className="w-28 px-4 py-3 text-right text-xs font-medium text-slate-500 uppercase tracking-wider">
            Actions
          </th>
        </tr>
      </thead>
      <tbody className="divide-y divide-slate-100">
        {facts.map(fact => (
          <FactTableRow
            key={fact.id}
            fact={fact}
            source={sources.find(s => s.id === fact.sourceId)}
            isSelected={selectedFacts.has(fact.id)}
            isExpanded={expandedFact === fact.id}
            onToggleSelect={() => onToggleSelect(fact.id)}
            onToggleExpand={() => onToggleExpand(fact.id)}
            onEdit={() => onEdit && onEdit(fact)}
            onDelete={() => onDelete && onDelete(fact)}
          />
        ))}
      </tbody>
    </table>
  );
};

const FactTableRow = ({ 
  fact, 
  source, 
  isSelected, 
  isExpanded,
  onToggleSelect, 
  onToggleExpand,
  onEdit, 
  onDelete 
}) => {
  const displayText = getFactDisplayText(fact);
  const answerText = getFactAnswerText(fact);
  const cards = mockCards[fact.id] || [];
  
  return (
    <>
      <tr className={`hover:bg-slate-50 transition-colors ${isSelected ? 'bg-sky-50' : ''}`}>
        {/* Checkbox */}
        <td className="w-12 px-4 py-3 align-top">
          <button onClick={onToggleSelect} className="p-1 hover:bg-slate-200 rounded transition-colors">
            {isSelected ? (
              <CheckSquare className="w-4 h-4 text-sky-500" />
            ) : (
              <Square className="w-4 h-4 text-slate-400" />
            )}
          </button>
        </td>
        
        {/* Content */}
        <td className="px-4 py-3">
          <button 
            onClick={onToggleExpand}
            className="w-full text-left group"
          >
            <div className="flex items-start gap-2">
              <ChevronRight className={`w-4 h-4 mt-0.5 text-slate-400 flex-shrink-0 transition-transform ${isExpanded ? 'rotate-90' : ''}`} />
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-slate-900 group-hover:text-sky-600 transition-colors line-clamp-1">
                  {displayText}
                </p>
                {fact.factType === 'basic' && answerText && (
                  <p className="text-xs text-slate-500 line-clamp-1 mt-0.5">{answerText}</p>
                )}
                {fact.tags && fact.tags.length > 0 && (
                  <div className="flex items-center gap-1 mt-1 flex-wrap">
                    {fact.tags.slice(0, 3).map(tag => (
                      <span key={tag} className="px-1.5 py-0.5 bg-slate-100 rounded text-xs text-slate-500">
                        {tag}
                      </span>
                    ))}
                    {fact.tags.length > 3 && (
                      <span className="text-xs text-slate-400">+{fact.tags.length - 3}</span>
                    )}
                  </div>
                )}
              </div>
            </div>
          </button>
        </td>
        
        {/* Type */}
        <td className="w-24 px-4 py-3 text-center align-top">
          <FactTypeBadge type={fact.factType} />
        </td>
        
        {/* Cards */}
        <td className="w-20 px-4 py-3 text-center align-top">
          <span className="text-sm font-medium text-slate-700">{fact.cardCount}</span>
        </td>
        
        {/* Due */}
        <td className="w-20 px-4 py-3 text-center align-top">
          {fact.dueCount > 0 ? (
            <span className="inline-flex items-center justify-center w-6 h-6 bg-sky-100 text-sky-600 rounded-full text-xs font-medium">
              {fact.dueCount}
            </span>
          ) : (
            <span className="text-sm text-slate-400">—</span>
          )}
        </td>
        
        {/* Actions */}
        <td className="w-28 px-4 py-3 text-right align-top">
          <div className="flex items-center justify-end gap-1">
            <button 
              onClick={(e) => { e.stopPropagation(); onEdit(); }}
              className="p-1.5 hover:bg-slate-200 rounded-lg transition-colors"
              title="Edit"
            >
              <Edit2 className="w-4 h-4 text-slate-500" />
            </button>
            <button 
              onClick={(e) => { e.stopPropagation(); onDelete(); }}
              className="p-1.5 hover:bg-red-100 rounded-lg transition-colors"
              title="Delete"
            >
              <Trash2 className="w-4 h-4 text-slate-500 hover:text-red-500" />
            </button>
            <button className="p-1.5 hover:bg-slate-200 rounded-lg transition-colors">
              <MoreHorizontal className="w-4 h-4 text-slate-500" />
            </button>
          </div>
        </td>
      </tr>
      
      {/* Expanded Card Details */}
      {isExpanded && (
        <tr className="bg-slate-50">
          <td colSpan={6} className="px-4 py-4">
            <div className="ml-12">
              {/* Fact Preview */}
              <div className="bg-white rounded-xl border border-slate-200 p-4 mb-4">
                <div className="flex items-start justify-between mb-3">
                  <h4 className="text-sm font-medium text-slate-700">Fact Preview</h4>
                  {source && (
                    <div className="flex items-center gap-1.5 text-xs text-slate-500">
                      <SourceIcon type={source.type} className="w-3.5 h-3.5" />
                      <span>{source.name}</span>
                    </div>
                  )}
                </div>
                
                {fact.factType === 'basic' && (
                  <div className="space-y-3">
                    <div>
                      <p className="text-xs text-slate-500 mb-1">Front</p>
                      <p className="text-sm text-slate-900">{displayText}</p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-500 mb-1">Back</p>
                      <p className="text-sm text-slate-700">{answerText}</p>
                    </div>
                  </div>
                )}
                
                {fact.factType === 'cloze' && (
                  <div className="space-y-3">
                    <div>
                      <p className="text-xs text-slate-500 mb-1">Cloze Text</p>
                      <p className="text-sm text-slate-900">{answerText}</p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-500 mb-1">Deletions</p>
                      <div className="flex flex-wrap gap-2">
                        {fact.content.fields.find(f => f.name === 'text')?.value.match(/\{\{c\d+::([^}]+)\}\}/g)?.map((match, i) => {
                          const content = match.match(/\{\{c(\d+)::([^}]+)\}\}/);
                          return (
                            <span key={i} className="px-2 py-1 bg-purple-100 text-purple-700 rounded text-xs font-mono">
                              c{content[1]}: {content[2]}
                            </span>
                          );
                        })}
                      </div>
                    </div>
                  </div>
                )}
                
                {fact.factType === 'image_occlusion' && (
                  <div className="space-y-3">
                    <div>
                      <p className="text-xs text-slate-500 mb-1">Title</p>
                      <p className="text-sm text-slate-900">{displayText}</p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-500 mb-1">Masked Regions ({fact.cardCount})</p>
                      <div className="flex flex-wrap gap-2">
                        {fact.content.fields.find(f => f.name === 'masks')?.value.map(mask => (
                          <span key={mask.id} className="px-2 py-1 bg-amber-100 text-amber-700 rounded text-xs">
                            {mask.label}
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>
                )}
              </div>
              
              {/* Cards List */}
              <div>
                <h4 className="text-sm font-medium text-slate-700 mb-2">Cards ({cards.length})</h4>
                <div className="space-y-2">
                  {cards.map(card => (
                    <div 
                      key={card.id}
                      className="flex items-center justify-between p-3 bg-white rounded-lg border border-slate-200"
                    >
                      <div className="flex items-center gap-3">
                        <CardStateBadge state={card.state} step={card.step} />
                        {card.elementId && (
                          <span className="text-xs font-mono text-slate-500 bg-slate-100 px-1.5 py-0.5 rounded">
                            {card.elementId}
                          </span>
                        )}
                      </div>
                      <div className="flex items-center gap-4 text-xs text-slate-500">
                        {card.due && (
                          <span>Due: {new Date(card.due).toLocaleDateString()}</span>
                        )}
                        {card.interval && (
                          <span>Interval: {card.interval}</span>
                        )}
                        <span>Reps: {card.reps}</span>
                        {card.lapses > 0 && (
                          <span className="text-red-500">Lapses: {card.lapses}</span>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </td>
        </tr>
      )}
    </>
  );
};

// ============================================================================
// FACT GRID COMPONENT
// ============================================================================

const FactGrid = ({ facts, sources, selectedFacts, onToggleSelect, onEdit, onDelete }) => {
  return (
    <div className="p-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {facts.map(fact => (
        <FactGridCard
          key={fact.id}
          fact={fact}
          source={sources.find(s => s.id === fact.sourceId)}
          isSelected={selectedFacts.has(fact.id)}
          onToggleSelect={() => onToggleSelect(fact.id)}
          onEdit={() => onEdit && onEdit(fact)}
          onDelete={() => onDelete && onDelete(fact)}
        />
      ))}
    </div>
  );
};

const FactGridCard = ({ fact, source, isSelected, onToggleSelect, onEdit, onDelete }) => {
  const displayText = getFactDisplayText(fact);
  const answerText = getFactAnswerText(fact);
  
  return (
    <div className={`bg-white rounded-xl border ${isSelected ? 'border-sky-300 ring-2 ring-sky-100' : 'border-slate-200'} overflow-hidden hover:shadow-md transition-all`}>
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-2 bg-slate-50 border-b border-slate-100">
        <div className="flex items-center gap-2">
          <button onClick={onToggleSelect} className="p-0.5 hover:bg-slate-200 rounded transition-colors">
            {isSelected ? (
              <CheckSquare className="w-4 h-4 text-sky-500" />
            ) : (
              <Square className="w-4 h-4 text-slate-400" />
            )}
          </button>
          <FactTypeBadge type={fact.factType} />
        </div>
        <div className="flex items-center gap-1">
          <button onClick={onEdit} className="p-1 hover:bg-slate-200 rounded transition-colors">
            <Edit2 className="w-3.5 h-3.5 text-slate-500" />
          </button>
          <button onClick={onDelete} className="p-1 hover:bg-red-100 rounded transition-colors">
            <Trash2 className="w-3.5 h-3.5 text-slate-500 hover:text-red-500" />
          </button>
        </div>
      </div>
      
      {/* Content */}
      <div className="p-4">
        <p className="text-sm font-medium text-slate-900 line-clamp-2 mb-2">{displayText}</p>
        {fact.factType === 'basic' && answerText && (
          <p className="text-xs text-slate-500 line-clamp-2">{answerText}</p>
        )}
        {fact.factType === 'cloze' && (
          <div className="flex flex-wrap gap-1 mt-2">
            {fact.content.fields.find(f => f.name === 'text')?.value.match(/\{\{c\d+/g)?.slice(0, 3).map((match, i) => (
              <span key={i} className="px-1.5 py-0.5 bg-purple-100 text-purple-600 rounded text-xs">
                {match.replace('{{', '')}
              </span>
            ))}
            {(fact.content.fields.find(f => f.name === 'text')?.value.match(/\{\{c\d+/g)?.length || 0) > 3 && (
              <span className="text-xs text-slate-400">+more</span>
            )}
          </div>
        )}
        {fact.factType === 'image_occlusion' && (
          <div className="flex items-center gap-1 mt-2">
            <Image className="w-3.5 h-3.5 text-amber-500" />
            <span className="text-xs text-slate-500">
              {fact.content.fields.find(f => f.name === 'masks')?.value.length || 0} regions
            </span>
          </div>
        )}
      </div>
      
      {/* Footer */}
      <div className="flex items-center justify-between px-4 py-2 bg-slate-50 border-t border-slate-100">
        <div className="flex items-center gap-3 text-xs text-slate-500">
          <span className="flex items-center gap-1">
            <Layers className="w-3.5 h-3.5" />
            {fact.cardCount}
          </span>
          {fact.dueCount > 0 && (
            <span className="flex items-center gap-1 text-sky-600 font-medium">
              <Clock className="w-3.5 h-3.5" />
              {fact.dueCount} due
            </span>
          )}
        </div>
        {source && (
          <div className="flex items-center gap-1 text-xs text-slate-400">
            <SourceIcon type={source.type} className="w-3 h-3" />
          </div>
        )}
      </div>
    </div>
  );
};

// ============================================================================
// CREATE FACT MODAL
// ============================================================================

const CreateFactModal = ({ isOpen, onClose, onSubmit }) => {
  const [factType, setFactType] = useState('basic');
  const [frontText, setFrontText] = useState('');
  const [backText, setBackText] = useState('');
  const [clozeText, setClozeText] = useState('');
  
  if (!isOpen) return null;
  
  const factTypes = [
    { id: 'basic', label: 'Basic', icon: Layers, description: 'Front and back card' },
    { id: 'cloze', label: 'Cloze Deletion', icon: Type, description: 'Fill in the blank' },
    { id: 'image_occlusion', label: 'Image Occlusion', icon: Image, description: 'Hide parts of an image' },
  ];
  
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-2xl bg-white rounded-2xl shadow-2xl overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-slate-200">
          <h2 className="text-lg font-semibold text-slate-900">Create New Fact</h2>
          <button onClick={onClose} className="p-2 hover:bg-slate-100 rounded-lg transition-colors">
            <X className="w-5 h-5 text-slate-500" />
          </button>
        </div>
        
        {/* Type Selector */}
        <div className="px-6 py-4 border-b border-slate-200">
          <p className="text-sm font-medium text-slate-700 mb-3">Fact Type</p>
          <div className="grid grid-cols-3 gap-3">
            {factTypes.map(type => (
              <button
                key={type.id}
                onClick={() => setFactType(type.id)}
                className={`p-4 rounded-xl border-2 text-left transition-all ${
                  factType === type.id
                    ? 'border-sky-500 bg-sky-50'
                    : 'border-slate-200 hover:border-slate-300'
                }`}
              >
                <type.icon className={`w-5 h-5 mb-2 ${factType === type.id ? 'text-sky-500' : 'text-slate-400'}`} />
                <p className={`font-medium text-sm ${factType === type.id ? 'text-sky-700' : 'text-slate-700'}`}>
                  {type.label}
                </p>
                <p className="text-xs text-slate-500 mt-0.5">{type.description}</p>
              </button>
            ))}
          </div>
        </div>
        
        {/* Content Fields */}
        <div className="px-6 py-4 max-h-[400px] overflow-y-auto">
          {factType === 'basic' && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Front</label>
                <textarea
                  value={frontText}
                  onChange={(e) => setFrontText(e.target.value)}
                  placeholder="Enter the question or prompt..."
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent resize-none"
                  rows={3}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Back</label>
                <textarea
                  value={backText}
                  onChange={(e) => setBackText(e.target.value)}
                  placeholder="Enter the answer..."
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent resize-none"
                  rows={3}
                />
              </div>
            </div>
          )}
          
          {factType === 'cloze' && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Cloze Text</label>
                <p className="text-xs text-slate-500 mb-2">
                  Wrap text in {'{{c1::text}}'} to create deletions. Use c1, c2, etc. for multiple cards.
                </p>
                <textarea
                  value={clozeText}
                  onChange={(e) => setClozeText(e.target.value)}
                  placeholder="{{c1::Mitochondria}} are the {{c2::powerhouse}} of the cell."
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl text-sm font-mono focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent resize-none"
                  rows={4}
                />
              </div>
              {clozeText && (
                <div className="p-3 bg-slate-50 rounded-lg">
                  <p className="text-xs font-medium text-slate-500 mb-2">Preview</p>
                  <p className="text-sm text-slate-700">
                    {clozeText.replace(/\{\{c\d+::([^}]+)\}\}/g, '[$1]')}
                  </p>
                </div>
              )}
            </div>
          )}
          
          {factType === 'image_occlusion' && (
            <div className="space-y-4">
              <div className="border-2 border-dashed border-slate-300 rounded-xl p-8 text-center">
                <Upload className="w-10 h-10 text-slate-400 mx-auto mb-3" />
                <p className="text-sm font-medium text-slate-600 mb-1">Upload an image</p>
                <p className="text-xs text-slate-500">PNG, JPG, or WebP up to 10MB</p>
                <button className="mt-4 px-4 py-2 bg-slate-100 hover:bg-slate-200 rounded-lg text-sm font-medium text-slate-600 transition-colors">
                  Choose File
                </button>
              </div>
              <p className="text-xs text-slate-500 text-center">
                After uploading, you'll be able to draw regions to hide on the image.
              </p>
            </div>
          )}
        </div>
        
        {/* Footer */}
        <div className="flex items-center justify-between px-6 py-4 border-t border-slate-200 bg-slate-50">
          <button 
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-slate-800 transition-colors"
          >
            Cancel
          </button>
          <div className="flex items-center gap-2">
            <button className="px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-200 rounded-lg transition-colors">
              Save & Create Another
            </button>
            <button 
              onClick={() => { onSubmit && onSubmit({ factType, frontText, backText, clozeText }); onClose(); }}
              className="px-4 py-2 bg-sky-500 hover:bg-sky-600 rounded-lg text-white text-sm font-medium transition-colors"
            >
              Create Fact
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

// ============================================================================
// REMAINING COMPONENTS (Simplified from original)
// ============================================================================

// Command Palette
const CommandPalette = ({ isOpen, onClose, notebooks, currentNotebook }) => {
  const [query, setQuery] = useState('');
  const inputRef = useRef(null);
  
  useEffect(() => {
    if (isOpen) {
      inputRef.current?.focus();
      setQuery('');
    }
  }, [isOpen]);
  
  if (!isOpen) return null;
  
  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[15vh]">
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-2xl bg-white rounded-2xl shadow-2xl overflow-hidden">
        <div className="flex items-center gap-3 px-4 py-4 border-b border-slate-200">
          <Search className="w-5 h-5 text-slate-400" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search facts, cards, sources..."
            className="flex-1 text-lg outline-none placeholder:text-slate-400"
          />
          <kbd className="px-2 py-1 bg-slate-100 rounded text-xs text-slate-500">esc</kbd>
        </div>
        <div className="p-4 text-center text-slate-500 text-sm">
          Start typing to search...
        </div>
      </div>
    </div>
  );
};

// Home Dashboard
const HomeDashboard = ({ notebooks, stats, onSelectNotebook, setView }) => {
  const totalDue = notebooks.reduce((sum, nb) => sum + nb.dueCount, 0);
  
  return (
    <div className="flex-1 overflow-y-auto bg-slate-50">
      <div className="max-w-6xl mx-auto p-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900">Welcome back!</h1>
          <p className="text-slate-500 mt-1">Here's your learning overview</p>
        </div>
        
        <div className="bg-gradient-to-br from-sky-500 to-cyan-600 rounded-2xl p-6 text-white mb-8">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sky-100 mb-1">Cards due today</p>
              <p className="text-4xl font-bold">{totalDue}</p>
              <p className="text-sky-100 mt-2 text-sm">{stats.cardsReviewedToday} reviewed so far</p>
            </div>
            <div className="flex items-center gap-6">
              <div className="text-center">
                <div className="flex items-center gap-1 justify-center">
                  <span className="text-lg">🔥</span>
                  <span className="text-2xl font-bold">{stats.currentStreak}</span>
                </div>
                <p className="text-sky-100 text-sm">Day streak</p>
              </div>
              <div className="relative">
                <ProgressRing progress={stats.averageRetention} size={80} stroke={6} trackColor="rgba(255,255,255,0.2)" progressColor="#ffffff" />
                <div className="absolute inset-0 flex items-center justify-center">
                  <span className="text-lg font-bold">{stats.averageRetention}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div>
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold text-slate-900">Your Notebooks</h3>
          </div>
          <div className="grid grid-cols-2 gap-4">
            {notebooks.map(nb => (
              <button
                key={nb.id}
                onClick={() => { onSelectNotebook(nb); setView('notebook'); }}
                className="bg-white hover:bg-slate-50 rounded-2xl p-5 text-left transition-colors border border-slate-200 hover:border-sky-200"
              >
                <div className="flex items-center gap-3 mb-3">
                  <span className="text-2xl">{nb.emoji}</span>
                  <div>
                    <h4 className="font-semibold text-slate-900">{nb.name}</h4>
                    <p className="text-sm text-slate-500">{nb.totalFacts} facts • {nb.totalCards} cards</p>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4 text-sm">
                    {nb.dueCount > 0 && (
                      <span className="text-sky-600 font-medium">{nb.dueCount} due</span>
                    )}
                    <span className="text-slate-500">{nb.retention}% retention</span>
                  </div>
                  <div className="flex items-center gap-1 text-slate-400">
                    <span className="text-orange-400">🔥</span>
                    <span className="text-sm">{nb.streak}</span>
                  </div>
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

// Notebooks Dropdown
const NotebooksDropdown = ({ notebooks, current, setCurrent, isOpen, setIsOpen, setView, isInNotebook }) => {
  const totalDue = notebooks.reduce((sum, nb) => sum + nb.dueCount, 0);
  
  return (
    <div className="relative">
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${isOpen || isInNotebook ? 'bg-slate-100 text-slate-900' : 'text-slate-600 hover:bg-slate-100'}`}
      >
        {isInNotebook && current ? (
          <>
            <span className="text-base">{current.emoji}</span>
            <span>{current.name}</span>
          </>
        ) : (
          <>
            <Book className="w-4 h-4" />
            <span>Notebooks</span>
          </>
        )}
        {totalDue > 0 && !isInNotebook && (
          <span className="bg-sky-100 text-sky-700 px-1.5 py-0.5 rounded text-xs font-medium">{totalDue}</span>
        )}
        <ChevronDown className={`w-4 h-4 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </button>
      
      {isOpen && (
        <>
          <div className="fixed inset-0 z-40" onClick={() => setIsOpen(false)} />
          <div className="absolute top-full left-0 mt-2 w-72 bg-white rounded-xl shadow-xl border border-slate-200 overflow-hidden z-50">
            <div className="p-2">
              {notebooks.map(nb => (
                <button
                  key={nb.id}
                  onClick={() => { setCurrent(nb); setIsOpen(false); setView('notebook'); }}
                  className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all ${nb.id === current?.id && isInNotebook ? 'bg-sky-50 text-sky-900' : 'hover:bg-slate-50'}`}
                >
                  <span className="text-xl">{nb.emoji}</span>
                  <div className="flex-1 text-left">
                    <div className="font-medium text-slate-900">{nb.name}</div>
                    <div className="text-xs text-slate-500">{nb.totalFacts} facts • {nb.streak} day streak</div>
                  </div>
                  {nb.dueCount > 0 && (
                    <span className="bg-sky-500 text-white text-xs px-2 py-0.5 rounded-full font-medium">{nb.dueCount}</span>
                  )}
                </button>
              ))}
            </div>
            <div className="border-t border-slate-100 p-2">
              <button className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg hover:bg-slate-50 text-slate-600">
                <div className="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
                  <Plus className="w-4 h-4" />
                </div>
                <span>Create new notebook</span>
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

// Review Launcher
const ReviewLauncher = ({ notebooks, onStartReview }) => {
  const [open, setOpen] = useState(false);
  const totalDue = notebooks.reduce((sum, nb) => sum + nb.dueCount, 0);
  
  return (
    <div className="relative">
      <button onClick={() => setOpen(!open)} className="flex items-center gap-2 px-3 py-2 bg-sky-500 hover:bg-sky-600 rounded-lg text-white text-sm font-medium transition-colors">
        <Zap className="w-4 h-4" />
        <span>Review</span>
        <span className="bg-white/20 px-1.5 py-0.5 rounded text-xs">{totalDue}</span>
      </button>
      
      {open && (
        <>
          <div className="fixed inset-0 z-40" onClick={() => setOpen(false)} />
          <div className="absolute top-full right-0 mt-2 w-72 bg-white rounded-xl shadow-xl border border-slate-200 overflow-hidden z-50">
            <div className="p-2 border-b border-slate-100">
              <button onClick={() => { onStartReview(null); setOpen(false); }} className="w-full flex items-center gap-3 p-2.5 rounded-lg bg-gradient-to-r from-sky-500 to-cyan-500 text-white hover:from-sky-600 hover:to-cyan-600 transition-all">
                <div className="w-8 h-8 bg-white/20 rounded-lg flex items-center justify-center"><Zap className="w-4 h-4" /></div>
                <div className="flex-1 text-left">
                  <div className="font-medium text-sm">Review All</div>
                  <div className="text-xs text-white/80">{totalDue} cards due</div>
                </div>
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
            <div className="p-2">
              {notebooks.filter(nb => nb.dueCount > 0).map(nb => (
                <button key={nb.id} onClick={() => { onStartReview(nb); setOpen(false); }} className="w-full flex items-center gap-2.5 px-2 py-2 rounded-lg hover:bg-slate-50 transition-all">
                  <span className="text-lg">{nb.emoji}</span>
                  <div className="flex-1 text-left text-sm font-medium text-slate-900">{nb.name}</div>
                  <span className="text-sm text-sky-600 font-medium">{nb.dueCount}</span>
                </button>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
};

// Notebook Sidebar
const NotebookSidebar = ({ notebook, sources, facts, isCollapsed, setIsCollapsed }) => {
  const [sourcesOpen, setSourcesOpen] = useState(true);
  const [factsOpen, setFactsOpen] = useState(true);
  
  const totalCards = facts.reduce((sum, f) => sum + f.cardCount, 0);
  const totalDue = facts.reduce((sum, f) => sum + f.dueCount, 0);
  
  if (isCollapsed) {
    return (
      <aside className="w-12 bg-slate-50 border-r border-slate-200 flex flex-col">
        <div className="p-2">
          <button onClick={() => setIsCollapsed(false)} className="w-8 h-8 hover:bg-slate-200 rounded-lg flex items-center justify-center transition-colors">
            <ChevronsRight className="w-4 h-4 text-slate-500" />
          </button>
        </div>
      </aside>
    );
  }
  
  return (
    <aside className="w-64 bg-slate-50 border-r border-slate-200 flex flex-col">
      <div className="px-3 py-3 border-b border-slate-200 flex items-center justify-between">
        <div className="flex items-center gap-2 min-w-0">
          <span className="text-lg">{notebook.emoji}</span>
          <h2 className="font-semibold text-slate-900 text-sm truncate">{notebook.name}</h2>
        </div>
        <div className="flex items-center gap-1">
          <button className="p-1.5 hover:bg-slate-200 rounded-lg transition-colors">
            <Settings className="w-4 h-4 text-slate-400" />
          </button>
          <button onClick={() => setIsCollapsed(true)} className="p-1.5 hover:bg-slate-200 rounded-lg transition-colors">
            <ChevronsLeft className="w-4 h-4 text-slate-400" />
          </button>
        </div>
      </div>
      
      <div className="flex-1 overflow-y-auto">
        {/* Quick Stats */}
        <div className="px-3 py-3 border-b border-slate-200">
          <div className="grid grid-cols-2 gap-2">
            <div className="bg-white rounded-lg p-2 border border-slate-200">
              <p className="text-xs text-slate-500">Facts</p>
              <p className="text-lg font-semibold text-slate-900">{facts.length}</p>
            </div>
            <div className="bg-white rounded-lg p-2 border border-slate-200">
              <p className="text-xs text-slate-500">Cards</p>
              <p className="text-lg font-semibold text-slate-900">{totalCards}</p>
            </div>
          </div>
        </div>
        
        {/* Sources Section */}
        <div className="border-b border-slate-200">
          <button onClick={() => setSourcesOpen(!sourcesOpen)} className="w-full flex items-center gap-2 px-3 py-2.5 hover:bg-slate-100 transition-colors">
            <ChevronRight className={`w-4 h-4 text-slate-400 transition-transform ${sourcesOpen ? 'rotate-90' : ''}`} />
            <FolderOpen className="w-4 h-4 text-slate-500" />
            <span className="flex-1 text-left text-sm font-medium text-slate-700">Sources</span>
            <span className="text-xs text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded">{sources.length}</span>
          </button>
          {sourcesOpen && (
            <div className="pb-2 px-2 space-y-1">
              {sources.map(source => (
                <button key={source.id} className="w-full flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-slate-100 text-left transition-colors">
                  <SourceIcon type={source.type} className="w-4 h-4 text-slate-500" />
                  <span className="flex-1 text-sm text-slate-700 truncate">{source.name}</span>
                </button>
              ))}
              <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-400 hover:text-sky-500 hover:bg-slate-100 rounded-lg transition-colors">
                <Upload className="w-4 h-4" />
                <span>Add source</span>
              </button>
            </div>
          )}
        </div>
        
        {/* Facts Section */}
        <div>
          <button onClick={() => setFactsOpen(!factsOpen)} className="w-full flex items-center gap-2 px-3 py-2.5 hover:bg-slate-100 transition-colors">
            <ChevronRight className={`w-4 h-4 text-slate-400 transition-transform ${factsOpen ? 'rotate-90' : ''}`} />
            <Layers className="w-4 h-4 text-slate-500" />
            <span className="flex-1 text-left text-sm font-medium text-slate-700">Facts</span>
            <span className="text-xs text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded">{facts.length}</span>
          </button>
          {factsOpen && (
            <div className="pb-2 px-2 space-y-1">
              <button className="w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm text-slate-700 bg-sky-100 text-sky-700 font-medium">
                <span>All Facts</span>
                <span>{facts.length}</span>
              </button>
              <button className="w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm text-slate-700 hover:bg-slate-100 transition-colors">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-sky-500" />
                  <span>Due for review</span>
                </div>
                <span className="text-sky-600 font-medium">{totalDue}</span>
              </button>
              <div className="pt-2 border-t border-slate-200 mt-2 space-y-1">
                <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-500 hover:text-slate-700 hover:bg-slate-100 rounded-lg transition-colors">
                  <Tag className="w-4 h-4" />
                  <span>Manage tags</span>
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </aside>
  );
};

// ============================================================================
// MAIN APP
// ============================================================================

export default function LearningStudioV6() {
  const [currentView, setCurrentView] = useState('home');
  const [currentNotebook, setCurrentNotebook] = useState(mockNotebooks[0]);
  const [commandPaletteOpen, setCommandPaletteOpen] = useState(false);
  const [notebooksDropdownOpen, setNotebooksDropdownOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [createFactModalOpen, setCreateFactModalOpen] = useState(false);
  
  const sources = mockSources[currentNotebook.id] || [];
  const facts = mockFacts[currentNotebook.id] || [];
  const totalDue = mockNotebooks.reduce((sum, nb) => sum + nb.dueCount, 0);
  
  useEffect(() => {
    const handleKeyDown = (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        setCommandPaletteOpen(true);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);
  
  const handleStartReview = (scope) => {
    console.log('Start review:', scope);
  };
  
  return (
    <div className="h-screen flex flex-col bg-slate-100">
      {/* Top Nav Bar */}
      <header className="bg-white border-b border-slate-200 px-4 py-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1">
            <button className="w-8 h-8 bg-gradient-to-br from-sky-500 to-cyan-500 rounded-lg flex items-center justify-center mr-2">
              <Box className="w-5 h-5 text-white" />
            </button>
            
            <button 
              onClick={() => setCurrentView('home')}
              className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${currentView === 'home' ? 'bg-slate-100 text-slate-900' : 'text-slate-600 hover:bg-slate-100'}`}
            >
              <Home className="w-4 h-4" />
              <span>Home</span>
            </button>
            
            <NotebooksDropdown 
              notebooks={mockNotebooks}
              current={currentNotebook}
              setCurrent={setCurrentNotebook}
              isOpen={notebooksDropdownOpen}
              setIsOpen={setNotebooksDropdownOpen}
              setView={setCurrentView}
              isInNotebook={currentView === 'notebook'}
            />
          </div>
          
          <div className="flex items-center gap-3">
            <button 
              onClick={() => setCommandPaletteOpen(true)}
              className="flex items-center gap-2 px-3 py-2 bg-slate-100 hover:bg-slate-200 rounded-lg transition-colors"
            >
              <Search className="w-4 h-4 text-slate-400" />
              <span className="text-slate-500 text-sm">Search...</span>
              <kbd className="px-1.5 py-0.5 bg-white rounded border border-slate-200 text-xs text-slate-400">⌘K</kbd>
            </button>
            
            {currentView === 'notebook' && (
              <div className="flex items-center gap-2 px-3 py-1.5 bg-slate-100 rounded-lg text-sm">
                <span className="text-orange-500">🔥</span>
                <span className="text-slate-600 font-medium">{currentNotebook.streak}d</span>
                <div className="w-px h-4 bg-slate-300" />
                <span className="text-slate-600 font-medium">{currentNotebook.retention}%</span>
              </div>
            )}
            
            <ReviewLauncher 
              notebooks={mockNotebooks}
              onStartReview={handleStartReview}
            />
          </div>
        </div>
      </header>
      
      {/* Main Content */}
      {currentView === 'home' ? (
        <HomeDashboard 
          notebooks={mockNotebooks}
          stats={mockGlobalStats}
          onSelectNotebook={setCurrentNotebook}
          setView={setCurrentView}
        />
      ) : (
        <div className="flex-1 flex overflow-hidden">
          <NotebookSidebar
            notebook={currentNotebook}
            sources={sources}
            facts={facts}
            isCollapsed={sidebarCollapsed}
            setIsCollapsed={setSidebarCollapsed}
          />
          
          <FactBrowser
            facts={facts}
            sources={sources}
            notebook={currentNotebook}
            onCreateFact={() => setCreateFactModalOpen(true)}
            onEditFact={(fact) => console.log('Edit fact:', fact)}
            onDeleteFact={(fact) => console.log('Delete fact:', fact)}
            onStartReview={() => handleStartReview(currentNotebook)}
          />
        </div>
      )}
      
      {/* Modals */}
      <CommandPalette 
        isOpen={commandPaletteOpen}
        onClose={() => setCommandPaletteOpen(false)}
        notebooks={mockNotebooks}
        currentNotebook={currentView === 'notebook' ? currentNotebook : null}
      />
      
      <CreateFactModal
        isOpen={createFactModalOpen}
        onClose={() => setCreateFactModalOpen(false)}
        onSubmit={(data) => console.log('Create fact:', data)}
      />
    </div>
  );
}
