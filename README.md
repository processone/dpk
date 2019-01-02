# Data Portability Kit

DPK is a Data Portability Kit. This is the Swift Army Knife that let you take back control on your online data.

Thanks to GDPR, online providers now have to offer takeout features for your data. This is a great opportunity
to get back your data. It is thus now a good time to manage and possibly publish them using open tools and platform.

That said, each providers will give you access to your data in a different format.

The goal of DPK is to provide a unified tool to extract your data in a unified, ready-to-use format.

The project will create a directory structure that is directly usable with your post in Markdown format and Metadata
in a consistent and unique JSON format.

With DPK, you get the change to get your data in a format that can really be reused, without any hidden dependencies.

## General principles

The goal of the tool is to produce content that is self contained. As such, we do not want to rely on third-party
services that can disappear at anytime. That's why as much as we can, we try to resolve short URLs to their final
target.

We also do not want to promote trackers. When using Twitter oembed for example, we sanitize the provided HTML and
thus we do not includes the `widget.js` Javascript tags.

## Shortlinks

URL Shorteners were popular when it was needed to share long links on Twitter, due to Tweet size limitations. Now, they
are mostly used for click on shared links. Short URL also hide the real link and if the short URL service disappear or
decide to redirect to another target, the original content will be lost.

That's why the toolkit provide methods to resolve short URLs and replace the short URL link with it's longer form. It
helps preserving the web link feature by removing middlemen.

## Twitter

You can ask Twitter to download your archive here: [Your Twitter Data](https://twitter.com/settings/your_twitter_data).  
You will receive a link to download your archive when ready.

When you got it, unzip the file and convert your tweets to a Markdown directory structure with the command:

```bash
go run cmd/twitter-to-md/twitter-to-md.go ~/Downloads/twitter-2018-12-27-abcd121212 posts
``` 

It will create a directory with your data in a format you can reuse with your blogging tool platform.

In the process, it will also embed a local representation of quoted tweets and replace shortened links with their
original value.
