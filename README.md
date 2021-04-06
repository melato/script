Go package to facilitate running external programs almost as easily as from a shell.

Example:
```
	s := script.Script{Trace: true}

	s.Run("ls", "-s")

	s.Cmd("ls").ToFile("ls.out")

	s.Cmd("cat").InputString("hello").Run()

	s.Cmd("cat").InputBytes([]byte("bytes")).Run()

	s.Cmd("cat").InputFile("ls.out").Run()

	pwd := s.Cmd("pwd").ToString()

	s.Cmd("ls", "-1").PipeTo("sort", "-r").Run()
  
	return s.Error()

```
