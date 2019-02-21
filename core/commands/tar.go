package commands

import (
	"fmt"
	"io"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	tar "github.com/ipfs/go-ipfs/tar"

	"gx/ipfs/QmQtQrtNioesAWtrx8csBvfY37gTe94d6wQ3VikZUjxD39/go-ipfs-cmds"
	"gx/ipfs/QmR3bNAtBoTN6xZ2HQNqpRQARcDoazH9jU6zKUNjFyQKWS/go-path"
	dag "gx/ipfs/QmScf5hnTEK8fDpRJAbcdMnKXpKUp1ytdymzXUbXDCFssp/go-merkledag"
	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
)

var TarCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Utility functions for tar files in ipfs.",
	},

	Subcommands: map[string]*cmds.Command{
		"add": tarAddCmd,
		"cat": tarCatCmd,
	},
}

var tarAddCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Import a tar file into ipfs.",
		ShortDescription: `
'ipfs tar add' will parse a tar file and create a merkledag structure to
represent it.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.FileArg("file", true, false, "Tar file to add.").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		enc, err := cmdenv.GetCidEncoder(req)
		if err != nil {
			return err
		}

		it := req.Files.Entries()
		file, err := cmdenv.GetFileArg(it)
		if err != nil {
			return err
		}

		node, err := tar.ImportTar(req.Context, file, nd.DAG)
		if err != nil {
			return err
		}

		c := node.Cid()

//<<<<<<< HEAD
		return cmds.EmitOnce(res, &AddEvent{
			Name: it.Name(),
			Hash: enc.Encode(c),
//=======
//		nd.Provider.Provide(c)
//
//		fi.FileName()
//		return cmds.EmitOnce(res, &coreiface.AddEvent{
//			Name: fi.FileName(),
//			Hash: c.String(),
//>>>>>>> Add provider to ipfs and provide when adding/fetching
		})
	},
	Type: AddEvent{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *AddEvent) error {
			fmt.Fprintln(w, out.Hash)
			return nil
		}),
	},
}

var tarCatCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Export a tar file from IPFS.",
		ShortDescription: `
'ipfs tar cat' will export a tar file from a previously imported one in IPFS.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("path", true, false, "ipfs path of archive to export.").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		p, err := path.ParsePath(req.Arguments[0])
		if err != nil {
			return err
		}

		root, err := core.Resolve(req.Context, nd.Namesys, nd.Resolver, p)
		if err != nil {
			return err
		}

		rootpb, ok := root.(*dag.ProtoNode)
		if !ok {
			return dag.ErrNotProtobuf
		}

		r, err := tar.ExportTar(req.Context, rootpb, nd.DAG)
		if err != nil {
			return err
		}

		return res.Emit(r)
	},
}
