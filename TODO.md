## Frontend
 [x] Any buttons that are only icon based should ALWAYS have a tooltip with a clear, concise description of what the button does.
 [x] Use lucide-maximize2 and lucide-minimize2 for expand and collapse panel buttons
 [x] Update source tabs to be pill based instead of underline
 [ ] improve source tab styling
 [ ] main pane expand, auto collapses side panels when they should remember their state after main pane is minimized

 - Source detail panel needs to be refactored


**NOTE:**
 - need generic capability of 'select text' -> take action (create cards, add to notes, summarize, explain, etc)
 - need capability to grab rectangular portion of screen to create image occlusion cards (could be from pdf, notes, video, etc)


## Large Sources
 - ideally large pdfs can be supported
 - they should be broken up into chapters and 
 - TODO: sources will likely be fairly common across users at a school, so we need a way to avoid duplicate processing of large files, OCR, embedding, etc.  For example, all nursing students in the same class will have the same PDF or video of a lecture.  


### Source Detail Refactor
 - Currently have 4 tabs, source, cards, summary, and chat
 - This forces the user to toggle between tabs to view different information about the source
 - To improve this UX for the user, lets add a new right sidebar that can be collapsed and expanded and resized (like the left sidebar)
 - The source should always be visible in the left pane and in the right pane the user can view cards, summary, and chat linked to that source as a set of tabs
 - Don't implement this now, but keep in mind for the design - the left source pane should be able to send information to the right sidebar and vice versa
  - for example, user could click on a card in the right sidebar and the left source pane should scroll to the section of the source that the card is linked to
  - or the as the user scrolls through the source in the left pane, the cards tab in the right sidebar should update to show cards linked to the current section of the source
  - the key is that information may need to flow between these 2 components, so bake this into the design from the start
 - As always, use best practices for svelte and component design.  Keep an eye on correctness, performance and simplicity.




### Source Viewers / Content Components
**PDF and Document Based Sources**
 - PDF displayed in the left pane
 - the right pane, would show cards, summary, and chat linked to that source
 - viewer allows for interactions to select text -> take action
 - viewer allows for interactions to select image -> take action
 - as user scrolls cards tab should update to show cards linked to that section of the pdf (later, more advanced features)

**Audio**
 - left pane has audio player with optional transcript below it
 - users can select text in transcript -> take action (should also capture timestamp if possible) 
 - as user scrolls cards tab should update to show cards linked to that section of the audio (later, more advanced features)

**Video**
 - left pane has video player with optional transcript below it
 - users can select text in transcript -> take action (should also capture timestamp if possible) 
 - as user watches video the cards tab should update to show cards linked to the current section of the video (later, more advanced features)



