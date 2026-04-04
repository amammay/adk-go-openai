# adk-go-openai

This repository provides a Go implementation of OpenAI support for ADK.
It is based on:

- https://github.com/levigross/adk-go/tree/2846ff82dc96c1ac1ecc0087633650c2e16a6546

## Why this repository exists

Unfortunately, this is a duplicate implementation, but it exists for practical reasons.

- https://github.com/google/adk-go/pull/242 was opened on Nov 9, 2025.
- A Googler commented on [Jan 28, 2026](https://github.com/google/adk-go/pull/242#issuecomment-3811798379) that a new
  `google/adk-go-community` repository would be set up to host this code.
- As of April 4, 2026, that repository has not been created.

In the meantime, multiple PRs and modules have attempted to solve the same problem.

### Related PRs

- https://github.com/google/adk-go/pull/242
- https://github.com/google/adk-go/pull/314
- https://github.com/google/adk-go/pull/342

### Related modules

- https://github.com/byebyebruce/adk-go-openai
- https://github.com/achetronic/adk-utils-go

After reviewing and trying those modules in my own projects, they did not behave exactly as needed.
The original PR from @levigross worked correctly for my use case, so I copied it here and shimmed in
the internal pieces from `adk-go` needed to make this repository work.

All credit goes to @levigross for the original implementation. Im hoping that someday this repo can be archived and a
more official implementation can be created.