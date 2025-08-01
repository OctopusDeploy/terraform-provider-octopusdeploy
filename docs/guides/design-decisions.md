---
page_title: "Design Decisions"
subcategory: "Guides"
---

# Design decisions

This document is to explain our stance on certain implementation decisions for the provider.

## Version-Controlled Projects and Terraform

When managing Deployment Processes and Runbooks, you can choose one of two options. Either:

1. Define a Database-backed Project and its Process(es) in Terraform

*or*

2. Define a Version-Controlled Project and its settings in Terraform, but define the Process as OCL in Git

~> We won’t allow you to try to manage the Process of a Version-Controlled Project in Terraform, and will give an informative error message if you try.

### Why
We believe if you're going to version control your config, then you should have one clear source-of-truth for each resource you're defining. Allowing you to define the process of a Version-Controlled Project through Terraform breaks that principle: which is the source of truth, the Terraform or the Project's OCL?

Terraform can’t nicely support all the richness that Version Control and Git provide (managing multiple branches, diverging processes etc), so attempting to allow you to define a Process in Terraform, and expect to have that write Git commits to the underlying OCL of a Version-Controlled project is a recipe for a bad time.

We want to let both Terraform, and Version Controlled Projects do what they do best, and not blur the lines between the two.

## Processes: Property Bags vs Strongly Typed Steps

For v1.0, we've moved away from strongly-typed steps, and leaned in to the key-value-pair Property bag.

### Why
In versions of the provider before v1.0, we experimented with both leaning on the Properties bag when defining a Step, and with providing strongly-typed Resources specific to each step. The result was a compromise between the two approaches: some properties were set via strongly typed attributes, and some needed to be set in the properties bag. 

It was hard to know which should be set where, and you could easily run into conflicts and state-drift when using a strongly-typed attribute when available (such as the `run_on_server` attribute), where Octopus Server would respond with elements in the property bag.

From a maintenance perspective, we currently don't have an easy way to generate these strongly-typed step representations, so each needs to be manually crafted and kept in sync with Octopus Server. This creates a maintenance burden that resulted in a bad experience for our Terraform Provider users. 

For the immediate future, the new Process resource will not have strongly-typed steps, and you'll need to use the properties bag instead. We've provided guidance on how to easily find the right properties for your step. 

In the medium-longer term, we would like to rearchitect the way our steps are defined in Octopus, which would enable us to lean on code generation to provide strongly typed step resources that are easily kept in sync with Octopus Server. The architectural improvements we've made in the Provider with the new Process resource set us up to be able to adopt strongly-typed steps without the previous state-drift issues, once we're able to do it in a maintainable way.
