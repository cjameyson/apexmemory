import React, { useState, useEffect, useRef } from 'react';
import { Book, FileText, ChevronDown, Plus, Search, Mic, Youtube, Link, StickyNote, Command, Star, ArrowRight, Check, Zap, X, Send, Upload, MessageSquare, ChevronRight, Eye, Box, Play, BookOpen, Grid, ChevronLeft, Tag, MoreHorizontal, Volume2, ExternalLink, Info, Edit2, Home, TrendingUp, Calendar, Target, Award, Clock, BarChart3, Layers, Mouse, Type, Image, Scissors, Edit, ArrowUpRight, Maximize2, Minimize2, SkipBack, SkipForward, Pause, RotateCcw, ZoomIn, ZoomOut, Download, Bookmark, List, RefreshCw, Menu, ChevronsLeft, ChevronsRight, Filter, Trash2, Settings, FolderOpen } from 'lucide-react';

// Mock Data
const mockNotebooks = [
  { id: 1, name: 'Biology 101', emoji: '🧬', color: '#0ea5e9', dueCount: 23, streak: 7, totalCards: 156, retention: 87 },
  { id: 2, name: 'Spanish B2', emoji: '🇪🇸', color: '#f59e0b', dueCount: 45, streak: 12, totalCards: 412, retention: 82 },
  { id: 3, name: 'Machine Learning', emoji: '🤖', color: '#6366f1', dueCount: 12, streak: 3, totalCards: 89, retention: 91 },
  { id: 4, name: 'History - WWII', emoji: '📜', color: '#ec4899', dueCount: 8, streak: 5, totalCards: 234, retention: 79 },
];

const mockSources = {
  1: [
    { id: 1, name: 'Chapter 3 - Cell Division', type: 'pdf', cards: 24, excerpt: 'Cell division is the process by which a parent cell divides into two or more daughter cells...', pages: 42, addedAt: '2 days ago' },
    { id: 2, name: 'Mitosis vs Meiosis', type: 'youtube', cards: 18, excerpt: 'A comprehensive comparison of the two types of cell division...', duration: '12:34', addedAt: '1 week ago' },
    { id: 3, name: 'Lab Notes - Microscopy', type: 'notes', cards: 8, excerpt: 'Observations from the microscopy lab session...', pages: 5, addedAt: '3 days ago' },
    { id: 4, name: 'Genetics Lecture Recording', type: 'audio', cards: 31, excerpt: 'Professor Smith discusses heredity and genetic expression...', duration: '48:22', addedAt: '1 hour ago' },
    { id: 5, name: 'Biology LibreTexts - Cells', type: 'url', cards: 15, excerpt: 'Open-source textbook covering cellular biology fundamentals...', addedAt: '5 days ago' },
  ],
};

const mockCards = {
  1: [
    { id: 1, front: 'What is the powerhouse of the cell?', back: 'The mitochondria - responsible for producing ATP through cellular respiration', sourceId: 1, due: true, interval: '1d', tags: ['organelles', 'exam-1'] },
    { id: 2, front: 'What are the phases of mitosis?', back: 'Prophase, Metaphase, Anaphase, Telophase (PMAT)', sourceId: 1, due: true, interval: '4h', tags: ['cell-division', 'exam-1'] },
    { id: 3, front: 'Define apoptosis', back: 'Programmed cell death - a controlled process of eliminating unwanted or damaged cells', sourceId: 2, due: true, interval: '2d', tags: ['cell-death'] },
    { id: 4, front: 'How many chromosomes do human somatic cells have?', back: '46 (23 pairs) - called the diploid number', sourceId: 3, due: false, interval: '1w', tags: ['genetics'] },
    { id: 5, front: 'What is crossing over?', back: 'Exchange of genetic material between homologous chromosomes during meiosis I', sourceId: 1, due: true, interval: '10m', tags: ['meiosis', 'genetics'] },
    { id: 6, front: 'What enzyme unwinds DNA during replication?', back: 'Helicase - it breaks hydrogen bonds between base pairs', sourceId: 4, due: true, interval: '1d', tags: ['dna', 'enzymes'] },
    { id: 7, front: 'Difference between chromatin and chromosomes?', back: 'Chromatin is loosely coiled DNA (during interphase); chromosomes are tightly condensed chromatin (during division)', sourceId: 1, due: false, interval: '2w', tags: ['cell-division'] },
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

// Utility Components
const SourceIcon = ({ type, className }) => {
  const icons = { pdf: <FileText className={className} />, youtube: <Youtube className={className} />, url: <Link className={className} />, audio: <Mic className={className} />, notes: <StickyNote className={className} /> };
  return icons[type] || <FileText className={className} />;
};

// Progress Ring
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

// Command Palette (Cmd+K)
const CommandPalette = ({ isOpen, onClose, notebooks, currentNotebook }) => {
  const [query, setQuery] = useState('');
  const [scope, setScope] = useState('notebook');
  const inputRef = useRef(null);
  
  useEffect(() => {
    if (isOpen) {
      inputRef.current?.focus();
      setQuery('');
    }
  }, [isOpen]);
  
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === 'Escape') onClose();
    };
    if (isOpen) window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);
  
  if (!isOpen) return null;
  
  const mockResults = [
    { type: 'card', title: 'What is the powerhouse of the cell?', subtitle: 'Biology 101 • Cell Division', icon: Layers },
    { type: 'card', title: 'Define apoptosis', subtitle: 'Biology 101 • Mitosis vs Meiosis', icon: Layers },
    { type: 'source', title: 'Chapter 3 - Cell Division', subtitle: 'Biology 101 • PDF • 24 cards', icon: FileText },
    { type: 'notebook', title: 'Machine Learning', subtitle: '89 cards • 12 due', icon: Book },
  ];
  
  const commands = [
    { label: 'Create new card', shortcut: 'C', icon: Plus },
    { label: 'Add source', shortcut: 'S', icon: Upload },
    { label: 'Start review', shortcut: 'R', icon: Zap },
    { label: 'Switch notebook', shortcut: 'N', icon: Book },
  ];
  
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
            placeholder={scope === 'global' ? 'Search all notebooks...' : `Search in ${currentNotebook?.name || 'notebook'}...`}
            className="flex-1 text-lg outline-none placeholder:text-slate-400"
          />
          <div className="flex items-center gap-1">
            <button
              onClick={() => setScope('notebook')}
              className={`px-2 py-1 rounded-md text-xs font-medium transition-colors ${scope === 'notebook' ? 'bg-sky-100 text-sky-700' : 'text-slate-500 hover:bg-slate-100'}`}
            >
              {currentNotebook?.emoji} This notebook
            </button>
            <button
              onClick={() => setScope('global')}
              className={`px-2 py-1 rounded-md text-xs font-medium transition-colors ${scope === 'global' ? 'bg-sky-100 text-sky-700' : 'text-slate-500 hover:bg-slate-100'}`}
            >
              🌐 All
            </button>
          </div>
        </div>
        
        <div className="max-h-[400px] overflow-y-auto">
          {query ? (
            <div className="p-2">
              <div className="px-3 py-2 text-xs font-medium text-slate-400 uppercase">Results</div>
              {mockResults.map((result, i) => (
                <button key={i} className="w-full flex items-center gap-3 px-3 py-3 rounded-xl hover:bg-slate-50 transition-colors text-left">
                  <div className="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center">
                    <result.icon className="w-5 h-5 text-slate-500" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-slate-900 truncate">{result.title}</div>
                    <div className="text-sm text-slate-500 truncate">{result.subtitle}</div>
                  </div>
                  <ArrowUpRight className="w-4 h-4 text-slate-300" />
                </button>
              ))}
            </div>
          ) : (
            <div className="p-2">
              <div className="px-3 py-2 text-xs font-medium text-slate-400 uppercase">Quick actions</div>
              {commands.map((cmd, i) => (
                <button key={i} className="w-full flex items-center gap-3 px-3 py-3 rounded-xl hover:bg-slate-50 transition-colors">
                  <div className="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center">
                    <cmd.icon className="w-5 h-5 text-slate-500" />
                  </div>
                  <span className="flex-1 text-left font-medium text-slate-700">{cmd.label}</span>
                  <kbd className="px-2 py-1 bg-slate-100 rounded text-xs text-slate-500 font-mono">⌘ {cmd.shortcut}</kbd>
                </button>
              ))}
            </div>
          )}
        </div>
        
        <div className="px-4 py-3 bg-slate-50 border-t border-slate-200 flex items-center justify-between text-xs text-slate-500">
          <div className="flex items-center gap-4">
            <span className="flex items-center gap-1"><kbd className="px-1.5 py-0.5 bg-white rounded border text-[10px]">↑↓</kbd> Navigate</span>
            <span className="flex items-center gap-1"><kbd className="px-1.5 py-0.5 bg-white rounded border text-[10px]">↵</kbd> Select</span>
            <span className="flex items-center gap-1"><kbd className="px-1.5 py-0.5 bg-white rounded border text-[10px]">esc</kbd> Close</span>
          </div>
          <span className="flex items-center gap-1"><Command className="w-3 h-3" /> <kbd className="px-1.5 py-0.5 bg-white rounded border text-[10px]">⌘K</kbd></span>
        </div>
      </div>
    </div>
  );
};

// Home Dashboard View
const HomeDashboard = ({ notebooks, stats, onSelectNotebook, setView }) => {
  const totalDue = notebooks.reduce((sum, nb) => sum + nb.dueCount, 0);
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
  
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
        
        <div className="grid grid-cols-4 gap-4 mb-8">
          {[
            { icon: Layers, label: 'Total Cards', value: stats.totalCards },
            { icon: Calendar, label: 'Due Tomorrow', value: stats.upcomingDue.tomorrow },
            { icon: Target, label: 'This Week', value: stats.upcomingDue.thisWeek },
            { icon: Award, label: 'Best Streak', value: `${stats.longestStreak} days` },
          ].map((stat, i) => (
            <div key={i} className="bg-white rounded-2xl p-5 border border-slate-200">
              <div className="flex items-center gap-2 text-slate-500 mb-2">
                <stat.icon className="w-4 h-4" />
                <span className="text-sm">{stat.label}</span>
              </div>
              <p className="text-2xl font-bold text-slate-900">{stat.value}</p>
            </div>
          ))}
        </div>
        
        <div className="mb-8">
          <h3 className="font-semibold text-slate-900 mb-4">This week</h3>
          <div className="bg-white rounded-2xl p-5 border border-slate-200">
            <div className="flex items-end justify-between h-32 gap-2">
              {stats.reviewsThisWeek.map((count, i) => (
                <div key={i} className="flex-1 flex flex-col items-center gap-2">
                  <div className="w-full bg-slate-100 rounded-lg overflow-hidden flex-1 flex items-end">
                    <div 
                      className={`w-full rounded-lg transition-all ${i < 5 ? 'bg-sky-500' : 'bg-slate-200'}`}
                      style={{ height: `${Math.max((count / 60) * 100, 4)}%` }}
                    />
                  </div>
                  <span className="text-xs text-slate-500">{days[i]}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
        
        <div>
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold text-slate-900">Your Notebooks</h3>
            <button className="text-sm text-sky-600 hover:text-sky-700 font-medium">View all</button>
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
                    <p className="text-sm text-slate-500">{nb.totalCards} cards</p>
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

// Notebooks Dropdown - Shows selected notebook or "Notebooks"
const NotebooksDropdown = ({ notebooks, current, setCurrent, isOpen, setIsOpen, setView, isInNotebook }) => {
  const totalDue = notebooks.reduce((sum, nb) => sum + nb.dueCount, 0);
  
  return (
    <div className="relative">
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${isOpen || isInNotebook ? 'bg-slate-100 text-slate-900' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}`}
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
                    <div className="text-xs text-slate-500">{nb.totalCards} cards • {nb.streak} day streak</div>
                  </div>
                  {nb.dueCount > 0 && (
                    <span className="bg-sky-500 text-white text-xs px-2 py-0.5 rounded-full font-medium">{nb.dueCount}</span>
                  )}
                </button>
              ))}
            </div>
            <div className="border-t border-slate-100 p-2">
              <button className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg hover:bg-slate-50 text-slate-600 transition-all">
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

// Collapsible Section Component
const SidebarSection = ({ title, icon: Icon, isOpen, onToggle, count, children, actions }) => {
  return (
    <div className="border-b border-slate-200 last:border-b-0">
      <button 
        onClick={onToggle}
        className="w-full flex items-center gap-2 px-3 py-2.5 hover:bg-slate-100 transition-colors"
      >
        <ChevronRight className={`w-4 h-4 text-slate-400 transition-transform ${isOpen ? 'rotate-90' : ''}`} />
        <Icon className="w-4 h-4 text-slate-500" />
        <span className="flex-1 text-left text-sm font-medium text-slate-700">{title}</span>
        {count !== undefined && (
          <span className="text-xs text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded">{count}</span>
        )}
        {actions && (
          <div className="flex items-center gap-1" onClick={(e) => e.stopPropagation()}>
            {actions}
          </div>
        )}
      </button>
      {isOpen && (
        <div className="pb-2">
          {children}
        </div>
      )}
    </div>
  );
};

// Source Type Specific Toolbar
const SourceToolbar = ({ source, onAction }) => {
  const toolbars = {
    pdf: (
      <div className="flex items-center gap-1">
        <div className="flex items-center gap-1 px-2 py-1 bg-slate-100 rounded-lg">
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Zoom out">
            <ZoomOut className="w-4 h-4 text-slate-600" />
          </button>
          <span className="text-sm text-slate-600 px-2">100%</span>
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Zoom in">
            <ZoomIn className="w-4 h-4 text-slate-600" />
          </button>
        </div>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <div className="flex items-center gap-1 px-2 py-1 bg-slate-100 rounded-lg">
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Previous page">
            <ChevronLeft className="w-4 h-4 text-slate-600" />
          </button>
          <span className="text-sm text-slate-600 px-2">1 / {source.pages}</span>
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Next page">
            <ChevronRight className="w-4 h-4 text-slate-600" />
          </button>
        </div>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Search in document">
          <Search className="w-4 h-4 text-slate-600" />
        </button>
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Table of contents">
          <List className="w-4 h-4 text-slate-600" />
        </button>
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Download">
          <Download className="w-4 h-4 text-slate-600" />
        </button>
      </div>
    ),
    youtube: (
      <div className="flex items-center gap-1">
        <div className="flex items-center gap-1 px-2 py-1 bg-slate-100 rounded-lg">
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Restart">
            <RotateCcw className="w-4 h-4 text-slate-600" />
          </button>
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Back 10s">
            <SkipBack className="w-4 h-4 text-slate-600" />
          </button>
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors bg-sky-100" title="Play">
            <Play className="w-4 h-4 text-sky-600" />
          </button>
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Forward 10s">
            <SkipForward className="w-4 h-4 text-slate-600" />
          </button>
        </div>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <span className="text-sm text-slate-600 px-2">0:00 / {source.duration}</span>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Loop section">
          <RefreshCw className="w-4 h-4 text-slate-600" />
        </button>
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Bookmark timestamp">
          <Bookmark className="w-4 h-4 text-slate-600" />
        </button>
        <button className="flex items-center gap-1 px-2 py-1 hover:bg-slate-100 rounded transition-colors text-sm text-slate-600">
          1x
        </button>
      </div>
    ),
    audio: (
      <div className="flex items-center gap-1">
        <div className="flex items-center gap-1 px-2 py-1 bg-slate-100 rounded-lg">
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Back 10s">
            <SkipBack className="w-4 h-4 text-slate-600" />
          </button>
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors bg-sky-100" title="Play">
            <Play className="w-4 h-4 text-sky-600" />
          </button>
          <button className="p-1.5 hover:bg-slate-200 rounded transition-colors" title="Forward 10s">
            <SkipForward className="w-4 h-4 text-slate-600" />
          </button>
        </div>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <span className="text-sm text-slate-600 px-2">16:08 / {source.duration}</span>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Volume">
          <Volume2 className="w-4 h-4 text-slate-600" />
        </button>
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Bookmark">
          <Bookmark className="w-4 h-4 text-slate-600" />
        </button>
        <button className="flex items-center gap-1 px-2 py-1 hover:bg-slate-100 rounded transition-colors text-sm text-slate-600">
          1x
        </button>
      </div>
    ),
    url: (
      <div className="flex items-center gap-1">
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Open original">
          <ExternalLink className="w-4 h-4 text-slate-600" />
        </button>
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Refresh">
          <RotateCcw className="w-4 h-4 text-slate-600" />
        </button>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Reader mode">
          <BookOpen className="w-4 h-4 text-slate-600" />
        </button>
      </div>
    ),
    notes: (
      <div className="flex items-center gap-1">
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Edit">
          <Edit2 className="w-4 h-4 text-slate-600" />
        </button>
        <div className="w-px h-6 bg-slate-200 mx-1" />
        <button className="p-1.5 hover:bg-slate-100 rounded transition-colors" title="Download">
          <Download className="w-4 h-4 text-slate-600" />
        </button>
      </div>
    ),
  };
  
  return (
    <div className="flex items-center justify-between px-4 py-2 bg-white border-b border-slate-200">
      {toolbars[source.type] || toolbars.notes}
      <div className="flex items-center gap-2">
        <button 
          onClick={() => onAction('generate')}
          className="flex items-center gap-1.5 px-3 py-1.5 bg-sky-50 hover:bg-sky-100 rounded-lg text-sky-700 text-sm font-medium transition-colors"
        >
          <Star className="w-4 h-4" />
          Generate Cards
        </button>
      </div>
    </div>
  );
};

// Selection Toolbar
const SelectionToolbar = ({ type, position, onAction, onClose }) => {
  if (!position) return null;
  
  const textActions = [
    { id: 'card', label: 'Create card', icon: Plus },
    { id: 'cloze', label: 'Cloze deletion', icon: Type },
    { id: 'highlight', label: 'Highlight', icon: Edit },
    { id: 'note', label: 'Add note', icon: StickyNote },
  ];
  
  const imageActions = [
    { id: 'occlusion', label: 'Image occlusion', icon: Image },
    { id: 'card', label: 'Card with image', icon: Plus },
  ];
  
  const audioActions = [
    { id: 'card', label: 'Card from segment', icon: Plus },
    { id: 'timestamp', label: 'Add timestamp note', icon: Clock },
  ];
  
  const actions = type === 'image' ? imageActions : type === 'audio' ? audioActions : textActions;
  
  return (
    <div 
      className="absolute z-50 bg-slate-900 rounded-xl shadow-2xl p-1 flex items-center gap-1"
      style={{ top: position.y, left: position.x, transform: 'translate(-50%, -100%) translateY(-8px)' }}
    >
      {actions.map(action => (
        <button
          key={action.id}
          onClick={() => { onAction(action.id); onClose(); }}
          className="flex items-center gap-2 px-3 py-2 rounded-lg text-white/80 hover:text-white hover:bg-white/10 transition-colors"
        >
          <action.icon className="w-4 h-4" />
          <span className="text-sm font-medium">{action.label}</span>
        </button>
      ))}
      <div className="w-px h-6 bg-white/20 mx-1" />
      <button onClick={onClose} className="p-2 rounded-lg text-white/50 hover:text-white hover:bg-white/10 transition-colors">
        <X className="w-4 h-4" />
      </button>
    </div>
  );
};

// Source Detail View
const SourceDetail = ({ source, cards, onClose, onStartReview, isExpanded, onToggleExpand }) => {
  const [activeTab, setActiveTab] = useState('source');
  const [selection, setSelection] = useState(null);
  const sourceCards = cards.filter(c => c.sourceId === source.id);
  const dueCards = sourceCards.filter(c => c.due);
  
  const handleMouseUp = (e, type = 'text') => {
    const selectedText = window.getSelection()?.toString();
    if (selectedText && selectedText.length > 3) {
      const rect = e.target.getBoundingClientRect();
      setSelection({
        type,
        position: { x: e.clientX - rect.left, y: e.clientY - rect.top },
        text: selectedText
      });
    }
  };
  
  const tabs = [
    { id: 'source', label: 'Source', icon: BookOpen },
    { id: 'cards', label: 'Cards', icon: Grid, count: sourceCards.length },
    { id: 'summary', label: 'Summary', icon: Info },
    { id: 'chat', label: 'Chat', icon: MessageSquare },
  ];
  
  return (
    <div className="flex-1 flex flex-col bg-white">
      <div className="px-4 py-3 border-b border-slate-200">
        <div className="flex items-center gap-3 mb-3">
          <button onClick={onClose} className="p-1.5 hover:bg-slate-100 rounded-lg transition-colors">
            <ChevronLeft className="w-5 h-5 text-slate-500" />
          </button>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 text-sm text-slate-500 mb-0.5">
              <SourceIcon type={source.type} className="w-4 h-4" />
              <span className="capitalize">{source.type}</span>
              {source.pages && <span>• {source.pages} pages</span>}
              {source.duration && <span>• {source.duration}</span>}
            </div>
            <h1 className="text-lg font-bold text-slate-900 truncate">{source.name}</h1>
          </div>
          <div className="flex items-center gap-2">
            {dueCards.length > 0 && (
              <button onClick={() => onStartReview(source)} className="flex items-center gap-2 px-3 py-1.5 bg-sky-500 hover:bg-sky-600 rounded-lg text-white text-sm font-medium transition-colors">
                <Play className="w-4 h-4" />Review ({dueCards.length})
              </button>
            )}
            <button onClick={onToggleExpand} className="p-1.5 hover:bg-slate-100 rounded-lg transition-colors" title={isExpanded ? 'Collapse' : 'Expand'}>
              {isExpanded ? <Minimize2 className="w-5 h-5 text-slate-500" /> : <Maximize2 className="w-5 h-5 text-slate-500" />}
            </button>
            <button className="p-1.5 hover:bg-slate-100 rounded-lg transition-colors">
              <MoreHorizontal className="w-5 h-5 text-slate-500" />
            </button>
          </div>
        </div>
        
        <div className="flex gap-1">
          {tabs.map(tab => (
            <button key={tab.id} onClick={() => setActiveTab(tab.id)} className={`flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${activeTab === tab.id ? 'bg-slate-100 text-slate-900' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-50'}`}>
              <tab.icon className="w-4 h-4" />
              {tab.label}
              {tab.count !== undefined && <span className={`px-1.5 py-0.5 rounded text-xs ${activeTab === tab.id ? 'bg-slate-200' : 'bg-slate-100'}`}>{tab.count}</span>}
            </button>
          ))}
        </div>
      </div>
      
      {activeTab === 'source' && (
        <SourceToolbar source={source} onAction={(action) => console.log('Toolbar action:', action)} />
      )}
      
      <div className="flex-1 overflow-y-auto">
        {activeTab === 'source' && (
          <div className="p-4">
            <div className="flex items-center gap-2 px-3 py-2 bg-sky-50 border border-sky-100 rounded-lg mb-4 text-sm text-sky-700">
              <Mouse className="w-4 h-4" />
              <span>Select text to create cards, highlights, or notes</span>
            </div>
            
            <div className="bg-slate-50 rounded-xl border border-slate-200 overflow-hidden relative">
              {source.type === 'pdf' && (
                <div className="p-6 min-h-[400px] relative" onMouseUp={(e) => handleMouseUp(e, 'text')}>
                  <div className="max-w-2xl mx-auto space-y-4 text-slate-700 leading-relaxed select-text">
                    <h2 className="text-xl font-bold text-slate-900">Chapter 3: Cell Division</h2>
                    <p>Cell division is the process by which a parent cell divides into two or more daughter cells. Cell division usually occurs as part of a larger cell cycle.</p>
                    <p><strong>Mitosis</strong> results in two daughter cells that are genetically identical to each other and to the original parent cell.</p>
                    <p><strong>Meiosis</strong>, on the other hand, produces four daughter cells, each with half the number of chromosomes of the parent cell.</p>
                  </div>
                  {selection && (
                    <SelectionToolbar type={selection.type} position={selection.position} onAction={(action) => console.log('Action:', action)} onClose={() => setSelection(null)} />
                  )}
                </div>
              )}
              
              {source.type === 'youtube' && (
                <div>
                  <div className="aspect-video bg-slate-900 flex items-center justify-center">
                    <div className="text-center">
                      <div className="w-16 h-16 bg-red-600 rounded-full flex items-center justify-center mx-auto mb-3 cursor-pointer hover:bg-red-700 transition-colors">
                        <Play className="w-7 h-7 text-white ml-1" />
                      </div>
                      <p className="text-white/60 text-sm">{source.duration}</p>
                    </div>
                  </div>
                  <div className="p-4" onMouseUp={(e) => handleMouseUp(e, 'text')}>
                    <h3 className="font-semibold text-slate-900 mb-2">Transcript</h3>
                    <div className="space-y-2 text-slate-600 select-text text-sm">
                      <p><span className="text-sky-600 font-mono text-xs">0:00</span> Welcome to today's video on mitosis and meiosis...</p>
                      <p><span className="text-sky-600 font-mono text-xs">1:24</span> Let's start by understanding what cell division actually means...</p>
                    </div>
                  </div>
                </div>
              )}
              
              {source.type === 'audio' && (
                <div className="p-6">
                  <div className="flex items-center gap-4 mb-4">
                    <div className="w-14 h-14 bg-sky-100 rounded-xl flex items-center justify-center">
                      <Mic className="w-7 h-7 text-sky-600" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-slate-900">{source.name}</h3>
                      <p className="text-slate-500 text-sm">{source.duration}</p>
                    </div>
                  </div>
                  <div className="h-20 bg-slate-100 rounded-lg mb-3 flex items-center justify-center">
                    <div className="flex items-end gap-0.5 h-12">
                      {Array.from({ length: 60 }).map((_, i) => (
                        <div key={i} className={`w-1 rounded-full ${i < 25 ? 'bg-sky-500' : 'bg-slate-300'}`} style={{ height: `${20 + Math.random() * 80}%` }} />
                      ))}
                    </div>
                  </div>
                </div>
              )}
              
              {(source.type === 'url' || source.type === 'notes') && (
                <div className="p-6" onMouseUp={(e) => handleMouseUp(e, 'text')}>
                  <div className="prose prose-slate max-w-none select-text">
                    <p className="text-slate-600 leading-relaxed">{source.excerpt}</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
        
        {activeTab === 'cards' && (
          <div className="p-4">
            <div className="flex items-center justify-between mb-3">
              <p className="text-sm text-slate-500">{sourceCards.length} cards from this source</p>
              <button className="flex items-center gap-1.5 px-2.5 py-1.5 bg-sky-50 hover:bg-sky-100 rounded-lg text-sky-700 text-sm font-medium transition-colors">
                <Plus className="w-4 h-4" />Create card
              </button>
            </div>
            <div className="space-y-2">
              {sourceCards.map(card => (
                <div key={card.id} className="p-3 bg-slate-50 rounded-xl border border-slate-100">
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex-1">
                      <p className="font-medium text-slate-900 text-sm">{card.front}</p>
                      <p className="text-sm text-slate-500 mt-1">{card.back}</p>
                    </div>
                    {card.due ? (
                      <span className="text-xs bg-sky-100 text-sky-600 px-2 py-0.5 rounded-full font-medium">Due</span>
                    ) : (
                      <span className="text-xs bg-emerald-100 text-emerald-600 px-2 py-0.5 rounded-full font-medium">{card.interval}</span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
        
        {activeTab === 'summary' && (
          <div className="p-4">
            <div className="bg-gradient-to-br from-sky-50 to-cyan-50 rounded-xl p-5 border border-sky-100">
              <div className="flex items-center gap-2 text-sm text-sky-600 mb-3">
                <Star className="w-4 h-4" />AI-generated summary
              </div>
              <h2 className="text-lg font-bold text-slate-900 mb-3">Cell Division Overview</h2>
              <p className="text-slate-600 text-sm leading-relaxed">Cell division is fundamental to life, enabling growth, repair, and reproduction.</p>
            </div>
          </div>
        )}
        
        {activeTab === 'chat' && (
          <div className="flex flex-col h-full">
            <div className="flex-1 p-4 flex flex-col items-center justify-center text-center">
              <div className="w-12 h-12 bg-sky-100 rounded-xl flex items-center justify-center mb-3">
                <MessageSquare className="w-6 h-6 text-sky-600" />
              </div>
              <h3 className="font-semibold text-slate-900 mb-1">Chat about this source</h3>
              <p className="text-slate-500 text-sm max-w-sm mb-4">Ask questions about "{source.name}"</p>
            </div>
            <div className="p-3 border-t border-slate-200">
              <div className="flex gap-2">
                <input type="text" placeholder="Ask about this source..." className="flex-1 px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent" />
                <button className="px-3 py-2 bg-sky-500 hover:bg-sky-600 text-white rounded-lg transition-colors">
                  <Send className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

// Cards Grid View
const CardsGridView = ({ cards, sources, viewMode, setViewMode }) => {
  const dueCards = cards.filter(c => c.due);
  const filteredCards = viewMode === 'due' ? dueCards : viewMode === 'mastered' ? cards.filter(c => !c.due) : cards;
  
  return (
    <div className="flex-1 flex flex-col bg-white">
      <div className="px-4 py-3 border-b border-slate-200">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-bold text-slate-900">All Cards</h2>
            <p className="text-sm text-slate-500">{cards.length} total • {dueCards.length} due</p>
          </div>
          <div className="flex items-center gap-2">
            <div className="flex bg-slate-100 rounded-lg p-0.5">
              {[{ id: 'all', label: 'All' }, { id: 'due', label: `Due (${dueCards.length})` }, { id: 'mastered', label: 'Mastered' }].map(v => (
                <button key={v.id} onClick={() => setViewMode(v.id)} className={`px-2.5 py-1 rounded-md text-sm font-medium transition-all ${viewMode === v.id ? 'bg-white shadow-sm text-slate-900' : 'text-slate-500 hover:text-slate-700'}`}>
                  {v.label}
                </button>
              ))}
            </div>
            <button className="p-1.5 hover:bg-slate-100 rounded-lg transition-colors">
              <Plus className="w-5 h-5 text-slate-500" />
            </button>
          </div>
        </div>
      </div>
      
      <div className="flex-1 overflow-y-auto p-4">
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-3">
          {filteredCards.map(card => {
            const source = sources.find(s => s.id === card.sourceId);
            return (
              <div key={card.id} className="bg-slate-50 rounded-xl border border-slate-100 overflow-hidden hover:border-slate-200 transition-colors">
                <div className="flex items-center justify-between px-3 py-1.5 bg-slate-100/50 border-b border-slate-100">
                  <div className="flex items-center gap-1.5 text-xs text-slate-500">
                    <SourceIcon type={source?.type} className="w-3 h-3" />
                    <span className="truncate max-w-[120px]">{source?.name}</span>
                  </div>
                  {card.due ? (
                    <span className="text-xs bg-sky-100 text-sky-600 px-1.5 py-0.5 rounded-full font-medium">Due</span>
                  ) : (
                    <span className="text-xs text-slate-400">{card.interval}</span>
                  )}
                </div>
                <div className="p-3">
                  <p className="font-medium text-slate-900 text-sm">{card.front}</p>
                  <p className="text-sm text-slate-500 mt-1.5 line-clamp-2">{card.back}</p>
                </div>
              </div>
            );
          })}
          <button className="p-4 rounded-xl border-2 border-dashed border-slate-200 text-slate-400 hover:border-sky-300 hover:text-sky-500 transition-all flex flex-col items-center justify-center gap-1.5 min-h-[100px]">
            <Star className="w-5 h-5" /><span className="text-sm font-medium">Create card</span>
          </button>
        </div>
      </div>
    </div>
  );
};

// Source List Item
const SourceListItem = ({ source, isSelected, onSelect }) => (
  <button onClick={onSelect} className={`w-full text-left px-3 py-2 rounded-lg transition-all ${isSelected ? 'bg-sky-100' : 'hover:bg-slate-100'}`}>
    <div className="flex items-center gap-2.5">
      <div className={`w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0 ${isSelected ? 'bg-sky-200' : 'bg-slate-100'}`}>
        <SourceIcon type={source.type} className={`w-4 h-4 ${isSelected ? 'text-sky-600' : 'text-slate-500'}`} />
      </div>
      <div className="flex-1 min-w-0">
        <h3 className={`font-medium text-sm truncate ${isSelected ? 'text-sky-900' : 'text-slate-900'}`}>{source.name}</h3>
        <div className="flex items-center gap-1.5 text-xs text-slate-500">
          <span>{source.cards} cards</span>
        </div>
      </div>
    </div>
  </button>
);

// Focus Mode
const FocusMode = ({ notebooks, selectedNotebook, selectedSource, onClose }) => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [flipped, setFlipped] = useState(false);
  const [completed, setCompleted] = useState([]);
  
  let allCards = [];
  let scopeLabel = '';
  
  if (selectedSource) {
    allCards = mockCards[selectedNotebook?.id || 1]?.filter(c => c.sourceId === selectedSource.id && c.due) || [];
    scopeLabel = selectedSource.name;
  } else if (selectedNotebook) {
    allCards = mockCards[selectedNotebook.id]?.filter(c => c.due) || [];
    scopeLabel = selectedNotebook.name;
  } else {
    notebooks.forEach(nb => {
      const nbCards = (mockCards[nb.id] || []).filter(c => c.due).map(c => ({ ...c, notebook: nb }));
      allCards.push(...nbCards);
    });
    scopeLabel = 'All Notebooks';
  }
  
  const currentCard = allCards[currentIndex];
  const sources = mockSources[selectedNotebook?.id || currentCard?.notebook?.id || 1] || [];
  const source = sources.find(s => s.id === currentCard?.sourceId);
  const progress = allCards.length > 0 ? (completed.length / allCards.length) * 100 : 100;
  
  const handleRate = () => {
    setCompleted([...completed, currentCard.id]);
    setTimeout(() => {
      if (currentIndex < allCards.length - 1) {
        setCurrentIndex(currentIndex + 1);
        setFlipped(false);
      }
    }, 200);
  };
  
  const isComplete = !currentCard || completed.length === allCards.length;
  
  return (
    <div className="fixed inset-0 bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 z-50 flex flex-col">
      <div className="flex items-center justify-between p-4">
        <button onClick={onClose} className="p-2 hover:bg-white/10 rounded-xl transition-colors">
          <X className="w-6 h-6 text-white/70" />
        </button>
        <div className="flex items-center gap-2 text-white/70">
          {currentCard?.notebook && <span className="text-xl">{currentCard.notebook.emoji}</span>}
          <span>{scopeLabel}</span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-white/70 text-sm">{completed.length}/{allCards.length}</span>
          <div className="w-32 h-2 bg-white/10 rounded-full overflow-hidden">
            <div className="h-full bg-gradient-to-r from-sky-400 to-cyan-400 rounded-full transition-all duration-300" style={{ width: `${progress}%` }} />
          </div>
        </div>
      </div>
      
      <div className="flex-1 flex items-center justify-center p-8">
        {isComplete ? (
          <div className="text-center">
            <div className="w-20 h-20 bg-emerald-500 rounded-full flex items-center justify-center mx-auto mb-6">
              <Check className="w-10 h-10 text-white" />
            </div>
            <h2 className="text-3xl font-bold text-white mb-2">All done!</h2>
            <p className="text-slate-400 mb-8">You've reviewed all due cards</p>
            <button onClick={onClose} className="px-6 py-3 bg-sky-500 hover:bg-sky-600 rounded-xl text-white font-medium transition-colors">Back to Studio</button>
          </div>
        ) : (
          <div className="w-full max-w-2xl">
            <div className="flex items-center justify-center gap-2 mb-4">
              <SourceIcon type={source?.type} className="w-4 h-4 text-white/40" />
              <span className="text-white/40 text-sm">{source?.name}</span>
            </div>
            
            <div className="bg-white rounded-3xl p-10 min-h-[300px] flex flex-col items-center justify-center cursor-pointer shadow-2xl" onClick={() => !flipped && setFlipped(true)}>
              {!flipped ? (
                <div className="text-center">
                  <p className="text-2xl text-slate-900 font-medium">{currentCard.front}</p>
                  <p className="text-sky-500 mt-6 flex items-center justify-center gap-2">
                    <Eye className="w-5 h-5" />Tap to reveal
                  </p>
                </div>
              ) : (
                <div className="text-center w-full">
                  <p className="text-slate-400 text-lg mb-4">{currentCard.front}</p>
                  <div className="w-24 h-px bg-slate-200 mx-auto my-6" />
                  <p className="text-2xl text-slate-900 font-medium">{currentCard.back}</p>
                </div>
              )}
            </div>
            
            {flipped && (
              <div className="flex gap-3 mt-8 justify-center">
                <button onClick={handleRate} className="px-8 py-4 bg-red-500/20 hover:bg-red-500/30 rounded-2xl text-red-300 font-medium transition-colors">
                  <div>Again</div>
                  <div className="text-xs opacity-70 mt-1">&lt;1m</div>
                </button>
                <button onClick={handleRate} className="px-8 py-4 bg-amber-500/20 hover:bg-amber-500/30 rounded-2xl text-amber-300 font-medium transition-colors">
                  <div>Hard</div>
                  <div className="text-xs opacity-70 mt-1">&lt;10m</div>
                </button>
                <button onClick={handleRate} className="px-8 py-4 bg-emerald-500/20 hover:bg-emerald-500/30 rounded-2xl text-emerald-300 font-medium transition-colors">
                  <div>Good</div>
                  <div className="text-xs opacity-70 mt-1">1d</div>
                </button>
                <button onClick={handleRate} className="px-8 py-4 bg-sky-500/20 hover:bg-sky-500/30 rounded-2xl text-sky-300 font-medium transition-colors">
                  <div>Easy</div>
                  <div className="text-xs opacity-70 mt-1">4d</div>
                </button>
              </div>
            )}
          </div>
        )}
      </div>
      
      <div className="p-4 flex justify-center">
        <div className="flex items-center gap-4 text-white/30 text-sm">
          <span>Space to flip</span><span>•</span><span>1-4 to rate</span><span>•</span><span>Esc to exit</span>
        </div>
      </div>
    </div>
  );
};

// Review Launcher
const ReviewLauncher = ({ notebooks, currentNotebook, onStartReview }) => {
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
              <div className="px-2 py-1.5 text-xs font-medium text-slate-400 uppercase">By notebook</div>
              {notebooks.filter(nb => nb.dueCount > 0).map(nb => (
                <button key={nb.id} onClick={() => { onStartReview(nb); setOpen(false); }} className="w-full flex items-center gap-2.5 px-2 py-2 rounded-lg transition-all hover:bg-slate-50">
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

// Notebook Sidebar Panel
const NotebookSidebar = ({ notebook, sources, cards, selectedSource, setSelectedSource, isCollapsed, setIsCollapsed }) => {
  const [sourcesOpen, setSourcesOpen] = useState(true);
  const [cardsOpen, setCardsOpen] = useState(true);
  
  const dueCards = cards.filter(c => c.due);
  const masteredCards = cards.filter(c => !c.due);
  
  if (isCollapsed) {
    return (
      <aside className="w-12 bg-slate-50 border-r border-slate-200 flex flex-col">
        <div className="p-2">
          <button 
            onClick={() => setIsCollapsed(false)}
            className="w-8 h-8 hover:bg-slate-200 rounded-lg flex items-center justify-center transition-colors"
            title="Expand sidebar"
          >
            <ChevronsRight className="w-4 h-4 text-slate-500" />
          </button>
        </div>
        <div className="flex-1 flex flex-col items-center py-2 gap-2">
          <button className="w-8 h-8 hover:bg-slate-200 rounded-lg flex items-center justify-center transition-colors" title="Sources">
            <FolderOpen className="w-4 h-4 text-slate-500" />
          </button>
          <button className="w-8 h-8 hover:bg-slate-200 rounded-lg flex items-center justify-center transition-colors" title="Cards">
            <Layers className="w-4 h-4 text-slate-500" />
          </button>
        </div>
      </aside>
    );
  }
  
  return (
    <aside className="w-72 bg-slate-50 border-r border-slate-200 flex flex-col">
      {/* Notebook Header */}
      <div className="px-3 py-3 border-b border-slate-200 flex items-center justify-between">
        <div className="flex items-center gap-2 min-w-0">
          <span className="text-lg">{notebook.emoji}</span>
          <h2 className="font-semibold text-slate-900 text-sm truncate">{notebook.name}</h2>
        </div>
        <div className="flex items-center gap-1">
          <button className="p-1.5 hover:bg-slate-200 rounded-lg transition-colors" title="Notebook settings">
            <Settings className="w-4 h-4 text-slate-400" />
          </button>
          <button 
            onClick={() => setIsCollapsed(true)}
            className="p-1.5 hover:bg-slate-200 rounded-lg transition-colors"
            title="Collapse sidebar"
          >
            <ChevronsLeft className="w-4 h-4 text-slate-400" />
          </button>
        </div>
      </div>
      
      {/* Scrollable Sections */}
      <div className="flex-1 overflow-y-auto">
        {/* Sources Section */}
        <SidebarSection
          title="Sources"
          icon={FolderOpen}
          isOpen={sourcesOpen}
          onToggle={() => setSourcesOpen(!sourcesOpen)}
          count={sources.length}
          actions={
            <button className="p-1 hover:bg-slate-200 rounded transition-colors" title="Add source">
              <Plus className="w-3.5 h-3.5 text-slate-400" />
            </button>
          }
        >
          <div className="px-2 space-y-1">
            <button 
              onClick={() => setSelectedSource(null)}
              className={`w-full text-left px-3 py-1.5 rounded-lg text-sm transition-all ${!selectedSource ? 'bg-sky-100 text-sky-700 font-medium' : 'text-slate-600 hover:bg-slate-100'}`}
            >
              All sources
            </button>
            {sources.map(source => (
              <SourceListItem
                key={source.id}
                source={source}
                isSelected={selectedSource?.id === source.id}
                onSelect={() => setSelectedSource(source)}
              />
            ))}
            <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-400 hover:text-sky-500 hover:bg-slate-100 rounded-lg transition-colors">
              <Upload className="w-4 h-4" />
              <span>Add source</span>
            </button>
          </div>
        </SidebarSection>
        
        {/* Cards Section */}
        <SidebarSection
          title="Cards"
          icon={Layers}
          isOpen={cardsOpen}
          onToggle={() => setCardsOpen(!cardsOpen)}
          count={cards.length}
          actions={
            <button className="p-1 hover:bg-slate-200 rounded transition-colors" title="Create card">
              <Plus className="w-3.5 h-3.5 text-slate-400" />
            </button>
          }
        >
          <div className="px-2 space-y-1">
            <button className="w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm text-slate-700 hover:bg-slate-100 transition-colors">
              <span>All cards</span>
              <span className="text-slate-400">{cards.length}</span>
            </button>
            <button className="w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm text-slate-700 hover:bg-slate-100 transition-colors">
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 rounded-full bg-sky-500" />
                <span>Due for review</span>
              </div>
              <span className="text-sky-600 font-medium">{dueCards.length}</span>
            </button>
            <button className="w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm text-slate-700 hover:bg-slate-100 transition-colors">
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 rounded-full bg-emerald-500" />
                <span>Mastered</span>
              </div>
              <span className="text-slate-400">{masteredCards.length}</span>
            </button>
            <div className="pt-2 border-t border-slate-200 mt-2">
              <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-500 hover:text-slate-700 hover:bg-slate-100 rounded-lg transition-colors">
                <Tag className="w-4 h-4" />
                <span>Manage tags</span>
              </button>
              <button className="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-500 hover:text-slate-700 hover:bg-slate-100 rounded-lg transition-colors">
                <Filter className="w-4 h-4" />
                <span>Filter cards</span>
              </button>
            </div>
          </div>
        </SidebarSection>
      </div>
    </aside>
  );
};

// Main App
export default function LearningStudioV5() {
  const [currentView, setCurrentView] = useState('home');
  const [currentNotebook, setCurrentNotebook] = useState(mockNotebooks[0]);
  const [selectedSource, setSelectedSource] = useState(null);
  const [viewMode, setViewMode] = useState('all');
  const [focusMode, setFocusMode] = useState(false);
  const [focusScope, setFocusScope] = useState({ notebook: null, source: null });
  const [commandPaletteOpen, setCommandPaletteOpen] = useState(false);
  const [notebooksDropdownOpen, setNotebooksDropdownOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [sourceExpanded, setSourceExpanded] = useState(false);
  
  const sources = mockSources[currentNotebook.id] || [];
  const cards = mockCards[currentNotebook.id] || [];
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
    if (scope === null) {
      setFocusScope({ notebook: null, source: null });
    } else if (scope.type) {
      setFocusScope({ notebook: currentNotebook, source: scope });
    } else {
      setFocusScope({ notebook: scope, source: null });
    }
    setFocusMode(true);
  };
  
  return (
    <div className="h-screen flex flex-col bg-slate-100">
      {/* Top Nav Bar */}
      <header className="bg-white border-b border-slate-200 px-4 py-2">
        <div className="flex items-center justify-between">
          {/* Left: Logo + Nav */}
          <div className="flex items-center gap-1">
            <button className="w-8 h-8 bg-gradient-to-br from-sky-500 to-cyan-500 rounded-lg flex items-center justify-center mr-2">
              <Box className="w-5 h-5 text-white" />
            </button>
            
            <button 
              onClick={() => { setCurrentView('home'); setSelectedSource(null); }}
              className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${currentView === 'home' ? 'bg-slate-100 text-slate-900' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}`}
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
          
          {/* Right: Search + Stats + Review */}
          <div className="flex items-center gap-3">
            {/* Search */}
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
              currentNotebook={currentNotebook}
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
          {/* Notebook Sidebar (collapsible with sections) */}
          {!sourceExpanded && (
            <NotebookSidebar
              notebook={currentNotebook}
              sources={sources}
              cards={cards}
              selectedSource={selectedSource}
              setSelectedSource={setSelectedSource}
              isCollapsed={sidebarCollapsed}
              setIsCollapsed={setSidebarCollapsed}
            />
          )}
          
          {/* Main Area */}
          {selectedSource ? (
            <SourceDetail 
              source={selectedSource}
              cards={cards}
              onClose={() => { setSelectedSource(null); setSourceExpanded(false); }}
              onStartReview={handleStartReview}
              isExpanded={sourceExpanded}
              onToggleExpand={() => setSourceExpanded(!sourceExpanded)}
            />
          ) : (
            <CardsGridView 
              cards={cards}
              sources={sources}
              viewMode={viewMode}
              setViewMode={setViewMode}
            />
          )}
        </div>
      )}
      
      {/* Command Palette */}
      <CommandPalette 
        isOpen={commandPaletteOpen}
        onClose={() => setCommandPaletteOpen(false)}
        notebooks={mockNotebooks}
        currentNotebook={currentView === 'notebook' ? currentNotebook : null}
      />
      
      {/* Focus Mode */}
      {focusMode && (
        <FocusMode 
          notebooks={mockNotebooks}
          selectedNotebook={focusScope.notebook}
          selectedSource={focusScope.source}
          onClose={() => setFocusMode(false)}
        />
      )}
    </div>
  );
}
