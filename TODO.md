A list of what I need to do, in order:

 - [x] Design UI
   - Probably straight in HTML
   - Won't do anything yet
   - Needs the API
 - [ ] Make API
   - Design routes. Maybe these:
     - `/query/[~, =, !]{query}`
     - `/query/[~, =, !]{query}/{limit | "first"}`
     - `/random`
     - `/random/{n}`
     - `/user/{id | name}`
     - `/eq/{id}`
	 - `/dump`
   - They return JSON
   - Implement them
 - [ ] Make web-server/-site
   - Frontend to API
   - Account creation
