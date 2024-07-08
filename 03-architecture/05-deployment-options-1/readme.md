# Deployment Options

We won't go in-depth with deployment instructions right now, that said, let's talk about how our choice of project architecture affects our deployment options, and how we _could_ deploy our application in the future. We'll only talk about cloud deployment options here, and by the "cloud" I'm just referring to a remote server that's managed by a third-party company like Google or Amazon.

![xkcd the cloud](https://imgs.xkcd.com/comics/the_cloud.png)

- [xkcd](https://xkcd.com/908/)

Using a cloud service to deploy applications is _super_ common these days because it's easy, fast, and cheap.

That said, it's still possible to deploy to a local or on-premise server, and some companies still do that, but it's not as common as it used to be.

## Monolithic deployment

Deploying a monolith is straightforward. Because your server is just one program, you just need to get it running on a server that's exposed to the internet and point your DNS records to it.

You could upload and run it on classic server, something like:

- AWS EC2
- GCP Compute Engine (GCE)
- Digital Ocean Droplets
- Azure Virtual Machines

Alternatively, you could use a platform that's specifically designed to run web applications, like:

- Heroku
- Google App Engine
- Fly.io
- AWS Elastic Beanstalk

## Decoupled deployment

With a decoupled architecture, you have _two_ different programs that need to be deployed. You would typically deploy your _back-end_ to the same kinds of places you would deploy a monolith.

For your front-end server, you can do the same, _or_ you can use a platform that's specifically designed to host static files and server-side rendered front-end apps, something like:

- Vercel
- Netlify
- GitHub Pages

Because the front-end bundle is likely just static files, you can host it easily on a [CDN (Content Delivery Network)](https://www.cloudflare.com/learning/cdn/what-is-a-cdn/) inexpensively.

## More powerful options

If you want to be able to scale your application up and down in specific ways, or you want to add other back-end servers to your stack, you might want to look into container orchestration options like Kubernetes and Docker Swarm.

## Don't worry about all this stuff!

I'm trying to gently introduce you to some popular technologies and how they work together, but you don't need to memorize all of these products and options.

# Quiz

Q: Multiple programs must be deployed in a **\_** architecture

A: Decoupled
