package gordon

type Repo struct {
	owner string
	name  string
}

func RepoApp(repo Repo) App {
	return App{
		owner: repo.owner,
		name:  repo.name,
	}
}

func (r Repo) Spec() AppSpec {
	return AppSpec{
		owner: r.owner,
		name:  r.name,
	}
}
