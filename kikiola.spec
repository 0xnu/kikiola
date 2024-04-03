Name: kikiola
Version: %{version}
Release: 1%{?dist}
Summary: Kikiola is a high-performance vector database written in Go.
License: MIT
URL: https://github.com/0xnu/kikiola
Source0: %{name}-%{version}.tar.gz

%description
Kikiola is a high-performance vector database written in Go.

%prep
%setup -q

%build
go build -o main ./cmd/main.go

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{_bindir}
cp main %{buildroot}%{_bindir}/%{name}

%clean
rm -rf %{buildroot}

%files
%{_bindir}/%{name}

%changelog
* %{date} Finbarrs Oketunji <f@finbarrs.eu>
- Initial package