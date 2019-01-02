# TODO

- Fix relative redirects (with Location: /something)
  For example:
   Processing link: https://donate.mozilla.org/
   => Resolved as /en-US/index.html
   Get /en-US/index.html: unsupported protocol scheme ""
- Render Youtube links
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
