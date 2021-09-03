#define FUSE_USE_VERSION 31

#include <fuse.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>
#include <stddef.h>
#include <assert.h>

static void *hello_init(struct fuse_conn_info *conn,
						struct fuse_config *cfg)
{
	(void)conn;
	cfg->kernel_cache = 1;
	//TODO allocate array and return it!
	return NULL;
}

static int arrayfs_getattr(const char *path, struct stat *stbuf,
						   struct fuse_file_info *fi)
{
	(void)fi;
	int res = 0;
	struct fuse_context fc = *fuse_get_context();
	fc.private_data;

	memset(stbuf, 0, sizeof(struct stat));
	if (strcmp(path, "/") == 0)
	{
		stbuf->st_mode = S_IFDIR | 0755;
		stbuf->st_nlink = 2; //TODO + num of elems on first level
	}
	else if (strchr(path + 1, "/") == NULL) // no more slashes in the path,
											// so it is a directory
	{
		stbuf->st_mode = S_IFDIR | 0755;
		stbuf->st_nlink = 2;
	}
	else if (strchr(path + 1, "/") == 1) // found second slash, so it is a file
	{
		stbuf->st_mode = S_IFREG | 0444;
		stbuf->st_nlink = 1;
		stbuf->st_size = strlen("hello"); //TODO length of the string contents
	}
	else
		res = -ENOENT;

	return res;
}

static int hello_readdir(const char *path, void *buf, fuse_fill_dir_t filler,
						 off_t offset, struct fuse_file_info *fi,
						 enum fuse_readdir_flags flags)
{
	(void)offset;
	(void)fi;
	(void)flags;

	if (strcmp(path, "/") != 0)
		return -ENOENT;

	filler(buf, ".", NULL, 0, 0);
	filler(buf, "..", NULL, 0, 0);
	filler(buf, "hello", NULL, 0, 0);

	return 0;
}

static int hello_open(const char *path, struct fuse_file_info *fi)
{
	if (strcmp(path + 1, "hello") != 0)
		return -ENOENT;

	if ((fi->flags & O_ACCMODE) != O_RDONLY)
		return -EACCES;

	return 0;
}

static int hello_read(const char *path, char *buf, size_t size, off_t offset,
					  struct fuse_file_info *fi)
{
	size_t len;
	(void)fi;
	if (strcmp(path + 1, "hello") != 0)
		return -ENOENT;

	len = strlen("hello");
	if (offset < len)
	{
		if (offset + size > len)
			size = len - offset;
		memcpy(buf, "hello" + offset, size);
	}
	else
		size = 0;

	return size;
}

static const struct fuse_operations hello_oper = {
	.init = hello_init,
	.getattr = arrayfs_getattr,
	.readdir = hello_readdir,
	.open = hello_open,
	.read = hello_read,
};

int main(int argc, char *argv[])
{
	int ret;
	struct fuse_args args = FUSE_ARGS_INIT(argc, argv);

	ret = fuse_main(args.argc, args.argv, &hello_oper, NULL);
	fuse_opt_free_args(&args);
	return ret;
}
