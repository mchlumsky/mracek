Name:           mracek
Version:        0.6.6
Release:        1%{?dist}
Summary:        mracek is a command line tool to manage your OpenStack configuration files

License:        MIT
URL:            https://github.com/mchlumsky/mracek
Source0:        https://github.com/mchlumsky/mracek/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang
BuildRequires:  git
Provides:       %{name} = %{version}

%description
mracek is a command line tool to manage your OpenStack configuration files.

%global debug_package %{nil}

%prep
%autosetup


%build
CGO_ENABLED=0 go build -v -o %{name}


%install
install -Dpm 0755 %{name} %{buildroot}/%{_bindir}/%{name}
install -Dpm 0644 completions/mracek.bash %{buildroot}/%{_datadir}/bash-completion/completions/%{name}
install -Dpm 0644 completions/mracek.zsh %{buildroot}/%{_datadir}/zsh/site-functions/_%{name}
install -Dpm 0644 completions/mracek.fish %{buildroot}/%{_datadir}/fish/vendor_completions.d/%{name}.fish


%files
/%{_bindir}/%{name}
/%{_datadir}/bash-completion/completions/%{name}
/%{_datadir}/zsh/site-functions/_%{name}
/%{_datadir}/fish/vendor_completions.d/%{name}.fish
%license LICENSE.md


%changelog
* Sat Oct 14 2023 Martin Chlumsky <martin.chlumsky@gmail.com> - 0.6.4-1
- First release
