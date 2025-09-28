- Start Date: 2025-09-17
- Target Version: 1.x

# Summary

Cancel Plan Preview Command is a command to cancel `plan-preview` commands that have not been picked up by piped to execute. Piped will continue to pull plan preview commands from control plane but will not execute them if there is a cancel plan preview command.

# Motivation

Currently, users can start plan preview commands but cannot cancel them. Therefore, we have the following two motivations to introduce a cancel plan preview command:

- Plan preview commands, especially terraform related plan preview commands, can be long running. This becomes a bottleneck especially when there are large number of commands but only a few number of piped agents to execute them. Cancelling the plan preview commands will shorten the queue of commands thereby saving resources and time.

- In the future we would like to have a generic way to cancel all types of commands. Exploring how to cancel plan preview command helps us developers to explore how to build a generic cancel command feature which could be extended to other commands.

# Detailed design

In this iteration of the cancel plan preview command, we limit ourselves to only trigger the cancel plan preview command using a github action. The following points describe how each component will handle the cancel plan preview command:

- Github Action: 
    - A new cancel-plan-preview github action will make a blocking call to pipectl to create a cancel plan preview command. This will contain the detail necessary to find the plan preview command to cancel, such as the head commit of the pull request.

- Pipectl: 
    - Makes a rpc tp create a cancel plan preview command. 
    - Periodically checks for the command result.

- Control Plane: 
    - Creates cancel plan preview commands for each piped configured to handle the repository and stores in the datastore.
    - Marks the command handled when Piped reports that the command was handled. 

- Piped:
    - Pulls and stores plan preview and cancel plan preview commands just like it pulls and stores any/every other command.
    - Before executing any plan preview command, Piped checks if there is a cancel plan preview command. If there is a cancel plan preview command then Piped does not execute the plan preview command. Piped reports to the control plane that the command was cancelled and sends the cancelled result to control plane.
    - If there is no cancel plan preview command, then the piped agent follows the process described in (Plan Preview RFC)[./0005-plan-preview.md]. 

The target version of piped agent for the cancel plan preview command is version 1. 

# Alternatives

- We rejected an idea to have a UI page to list all the currently running plan preview commands and let the users cancel the plan preview commands from there. The main motivation behind this idea was to address the difficulty in finding the PRs which are currently running the plan preview commands. This idea is not completely rejected, we want to first address how the plan preview is mostly used. We would possibly revisit this idea in the future in the context of a generic cancel command.

- We rejected an idea where the user would have to manually specify details of the plan preview command so that the cancel plan preview command could find the plan preview command to cancel it. We rejected this idea in favour of using github actions because github actions already provide us with this detail making it the most efficient method for users to trigger a cancel plan preview command.

# Unresolved questions

In this iteration, we do not cancel commands that piped is already executing. In a future iteration we would explore how to cancel commands that a piped is executing.
