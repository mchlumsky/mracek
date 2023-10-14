Name:           mracek
Version:        0.6.4
Release:        1%{?dist}
Summary:        mracek is a command line tool to manage your OpenStack configuration files

License:        MIT
URL:            https://github.com/mchlumsky/mracek
Source0:        https://github.com/mchlumsky/mracek/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang
BuildRequires:  git
Provides:       %{name} = %{version}

%description
mracek (Czech word meaning "little cloud") is a small command line tool to manage your OpenStack configuration files.

%global debug_package %{nil}

%prep
%autosetup


%build
CGO_ENABLED=0 go build -v -o %{name}


%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}
install -Dpm 0644 completions/mracek.bash %{buildroot}%{bash_completions_dir}/%{name}
install -Dpm 0644 completions/mracek.zsh %{buildroot}%{zsh_completions_dir}/_%{name}
install -Dpm 0644 completions/mracek.fish %{buildroot}%{fish_completions_dir}/%{name}.fish


%files
%{_bindir}/%{name}
%{bash_completions_dir}/%{name}
%{zsh_completions_dir}/_%{name}
%{fish_completions_dir}/%{name}.fish
%license LICENSE.md


%changelog
* Sat Oct 14 2023 Martin Chlumsky <martin.chlumsky@gmail.com> - 0.6.4-1
- First release
