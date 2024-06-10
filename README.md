# AI Fake News

A lightweight ai generated news site, for fun.

Note: It's clearly marked "fake news" to not cause confusion.

## About

Theres two types of pages you can view.

- The front page (via accessing "/" path)
  - Each day the front page is generated with headlines from the frontpage of https://bloomberg.com.
    - Some slight headline modifications might exist.
    - 100% ai generated articles based on the headline and certain prompts.
- Custom pages (the fun part)
  - Allows for generating custom pages on the fly, with custom prompts.
  - Requires a news category and prompt. Example: "Category: ["Technology"], Prompt: articles about technology ai apocolypse"
  - Custom pages are queued and processed. When the page is ready, it can be viewed on the path "/topics/{category}/{headline-article-name}"
  - Custom pages expire after 1 month unless marked "keep forever".


### Technology Used

1. Go for backend
    - html-web-crawler to fetch real articles
    - Gemini for ai generated content
2. Go for frontend!
    - html/template to build webpage without javascript
    - Simple fetch javascript used to generate new pages and link to them.
3. dynamodb to store information permantently
   
Total site in a 20MB binary, which can be run directly or from docker container. If you don't provide aws info, it will run in-memory without a database.
