# Which is Better?

There is _always_ a trade-off.

## Pros for monoliths

- Simpler to get started with
- Easier to deploy new versions because everything is always in sync
- In the case of the data being embedded in the HTML, the performance can result in better UX and SEO

## Pros for decoupled architectures

- Easier to scale as traffic grows
- Easier to practice good separation of concerns as the codebase grows
- Can be hosted on separate servers and using separate technologies
- Embedding data in the HTML is still possible with pre-rendering (similar to how Next.js works), it's just more complicated

## Can we have the best of both worlds?

Perhaps. My recommendation to someone building a new application from scratch would be to start with a monolith, but to keep the API and the front-end decoupled logically within the project from the start (like we're doing with Chirpy).

That way, our app is easy to get started with, but we can migrate to a fully decoupled architecture later if we need to.

[Monoliths vs decoupled front and back ends | engineer explains](https://www.youtube.com/watch?v=k0XuJjZ_HP0)

# Quiz

Q: Which will probably be cheaper to host as traffic and scope increases?

A: Decoupled
