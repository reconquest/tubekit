# TUBEKIT
[![All Contributors](https://img.shields.io/badge/all_contributors-3-orange.svg?style=flat-square)](#contributors)

tubekit helps you to operate Kubernetes clusters more effectively by reducing
time spent on everyday routines.

## tubectl

[![Build Status](https://travis-ci.org/reconquest/tubekit.svg?branch=master)](https://travis-ci.org/reconquest/tubekit)

Tubectl is a simple yet powerful wrapper around kubectl which adds a bit of magic
to your everyday kubectl routines by reducing the complexity of working with
contexts, namespaces and intelligent matching resources.

## Installation

#### Arch Linux

The package is always available in AUR:
[tubekit-git](https://aur.archlinux.org/packages/tubekit-git).

#### Debian/Fedora/CentOS

Packages can be found on [releases page](https://github.com/reconquest/tubekit/releases).

#### macOS

```
brew tap reconquest/tubekit https://github.com/reconquest/tubekit.git
brew install tubekit
```

#### Binaries

Binaries can be found on [releases page](https://github.com/reconquest/tubekit/releases).

#### Standalone

With your [GOPATH](https://github.com/golang/go/wiki/GOPATH) already set up:

```
go get github.com/reconquest/tubekit/cmd/...
```

## Usage

[![asciicast](https://asciinema.org/a/233185.svg)](https://asciinema.org/a/233185)

### Quick context

Very often you need to switch between environments and therefore you need to
switch your kubectl's context or specify it in kubectl's command line with the flag
`--context`. Well, there is even no short option for that flag. This is where
tubectl starts its magic.

You can just specify `@` followed by a name of the context, you don't need to specify
full name of the context, it will be matched by the existing list of contexts
specified in your kubectl configs.

Instead of typing `kubectl --context staging` or `kubectl --context production`
now you just type `tubectl @st` or `tubectl @prod`. It doesn't mean that
`--context` will no longer work because it will. `tubectl` does not break
`kubectl` as it has full backwards compatibility with the original kubectl.

### Quick namespaces

Boring to type `-n monitoring`, huh? Here is the solution — don't type so much,
just type `+mon` with tubectl. tubectl retrieves a list of namespaces from
a cluster and matches given namespace. It's fully back compatible with `-n
--namespace` flags.

Forget about `--all-namespaces`, just use the `++` flag.

With kubectl you used to do like that:

```
kubectl --context staging --all-namespaces get pods
```

Now you can type:
```
tubectl @st ++ get pods
```

Instead of `kube-system` in such command:
```
kubectl get pods -n kube-system
```

You can now type like this:
```
tubectl get pods +sys
```

Also, you can type `++` and `@st` at the end of cmdline or in the middle of
it, it really doesn't matter — tubectl will figure it out.

### Resource matching

That is the real panacea for several problems if you ever had a thought like:
* tailing logs of multiple containers at the same time
* executing a command in multiple containers in serial or parallel mode
* deleting/retrieving/patching all pods/deployments/ingresses matching given
    regular expression
* just want to attach to the first container of the pod with several replicas whatever
    suffix it might have

Meet the matching operator — `%`. 

For example, we have two replicas of nginx running in a cluster:
```
nginx-59d4c6bbcf-dft5z                  1/1     Running   8          18d
nginx-59d4c6bbcf-dvfd2                  1/1     Running   8          18d
```

Instead of getting the list of pods and then running `kubectl logs` against all
of them you can do the following:

```
tubectl logs nginx%
```

You need to put `%` at the end of your matching query, tubectl will figure out
what resource are you trying to match, will retrieve list of these resources
from given cluster with given namespace, will match it and then run your
command against these resources with all other arguments you've passed.

When you specify one `%` at the end of the query tubectl will run the command
against all matched resources. If you want to run only against a specific
replica, you need to specify its number like `%:X`, example:

```
tubectl logs nginx%:1
```

In this example, tubectl will retrieve all pods matching `nginx` and then get
1st pod from the list, the same can be applied to the second pod — `%:2`.

**What if you want to run the command against all the pods but in parallel mode?**

Similar to `++` for `--all-namespaces` tubectl introduces an operator `%%`,
let's say you want to tail logs of all nginx pods at the same time, then you
just do the following:

```
tubectl logs nginx%% -f
```

tubectl will match resources and will run `kubectl logs` against all pods at
the same time.

### Other advanced examples

* Deleting all deployments matching word `backend` in `staging` namespace

    ```
    tubectl delete deployment +stag backend%
    ```

    or in parallel (concurrent) mode:

    ```
    tubectl delete deployment +stag backend%%
    ```

* Describing all statefulsets matching word `redis` in `staging` namespace

    ```
    tubectl describe sts +stag redis%%
    ```

    or in parallel (concurrent) mode:

    ```
    tubectl describe sts +stag redis%%
    ```

* Executing a command in all pods matching word `apiserver` or `scheduler` in
    `kube-system` namespace 

    ```
    tubectl exec +sys '(apiserver|scheduler)%' -- id
    ```

    or in parallel (concurrent) mode:

    ```
    tubectl exec +sys '(apiserver|scheduler)%%' -- id
    ```

### Custom flags

Tubectl supports a few own flags, they all have prefix `--tube`:

* `--tube-version` - prints version of the program
* `--tube-debug` - enables debug mode, also can be turned on by `TUBEKIT_DEBUG`
* `--tube-help` - prints short help message about the program.

### Authors

* [Egor Kovetskiy](https://github.com/kovetskiy)
* [Stanislav Seletskiy](https://github.com/seletskiy)

Hire us, we do reduce business costs by optimizing the things.

### License

MIT

## Contributors

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore -->
<table><tr><td align="center"><a href="https://github.com/saromanov"><img src="https://avatars1.githubusercontent.com/u/7413968?v=4" width="100px;" alt="Sergey Romanov"/><br /><sub><b>Sergey Romanov</b></sub></a><br /><a href="https://github.com/reconquest/tubekit/commits?author=saromanov" title="Code">💻</a></td><td align="center"><a href="http://dmitry.shaposhnik.name"><img src="https://avatars1.githubusercontent.com/u/27363?v=4" width="100px;" alt="Dmitry Shaposhnik"/><br /><sub><b>Dmitry Shaposhnik</b></sub></a><br /><a href="#platform-e1senh0rn" title="Packaging/porting to new platform">📦</a></td><td align="center"><a href="https://twitter.com/nesl247"><img src="https://avatars3.githubusercontent.com/u/1037526?v=4" width="100px;" alt="Harrison Heck"/><br /><sub><b>Harrison Heck</b></sub></a><br /><a href="#platform-nesl247" title="Packaging/porting to new platform">📦</a></td></tr></table>

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
