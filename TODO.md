# TODO

- Fix relative redirects (with Location: /something)
  For example:
   Processing link: https://donate.mozilla.org/
   => Resolved as /en-US/index.html
   Get /en-US/index.html: unsupported protocol scheme ""
- Resolve HTML 5 / RDFa prefixes properly when parsing page.
- Render Youtube links
  It should work by embedding content in a way that avoid tracking. We cannot just embed Youtube video snippet.
  The archives converter should:
  - Download the video cover image locally
  - Generate a preview that do not leak information to Youtube.
  As a first step, clicking the video will get the browser to load the Youtube page, but later on, we can
  render the video locally by loading the content from Youtube inline. It will have the same effect as loading the
  Youtube page, but inline. The "on-demand" video play without forced Youtube script embed will be compliant with
  Do not track policy (Like what Medium is doing for example).
  See: https://webdesign.tutsplus.com/tutorials/how-to-lazy-load-embedded-youtube-videos--cms-26743
  http://coffeespace.org.uk/dnt.js 
- Resolve twitter short url inside embedded tweets.
- Generate entries for liked tweets ? They are not included in archive, so requires querying Twitter API to get them.
  We could just generate link.
- Add media types to metadata file
- Add metadata to SQLite index at root dir.
- Rename media file to shorter / more friendly filenames.
- Write initial test suite.
- Refactor / clean-up
- Fix JS header removal to make it more generic (in case there is a part1, etc.)
- Convert smileys to Emoji
- Use similarity to find duplicate post across several source of data
- Remove utm_ parameters from links (used for tracking promo campaigns)

## URL sanitization

- Fix m.engadget.com links (the domain does not exist anymore)
Example:
 http://m.engadget.com/default/article.do?artUrl=http://www.engadget.com/2011/02/08/nokia-ceo-stephen-elop-rallies-troops-in-brutally-honest-burnin/&category=classic&postPage=1
 => http://www.engadget.com/2011/02/08/nokia-ceo-stephen-elop-rallies-troops-in-brutally-honest-burnin/&category=classic&postPage=1

## Roadmap

Other possible services to support for archive cleaning and unification:

- Instagram
- Facebook
- Google+
- Hangout
- Medium
- LinkedIn
- Pinterest
- Flickr
- Quora
- Pocket
- Pinboard
- Dropbox Paper
- Evernote
- Feedly
