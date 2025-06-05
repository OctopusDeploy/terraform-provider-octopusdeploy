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
