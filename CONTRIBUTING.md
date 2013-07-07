## GitHub Workflow

Developing patches should follow this workflow:

### Initial Setup

1.  Fork on GitHub (click Fork button).  This creates your own working copy on github.
2.  Clone to computer: `git clone git@github.com:<<Your Github Username>>/speedtest.git`
3.  cd into your repo: `cd speedtest`
4.  Set up remote upstream: `git remote add -f upstream git://github.com/zpeters/speedtest.git`

### Adding a Feature

1.  Create a branch for the new feature: `git checkout -b my_new_feature`
2.  Work on your feature, add and commit as usual

Creating a branch is not strictly necessary, but it makes it easy to delete your branch when the feature has been merged into upstream, diff your branch with the version that actually ended in upstream, and to submit pull requests for multiple features (branches).

### Pushing to GitHub

8.  Push branch to GitHub: `git push origin my_new_feature`
9.  Issue pull request: Click Pull Request button on GitHub