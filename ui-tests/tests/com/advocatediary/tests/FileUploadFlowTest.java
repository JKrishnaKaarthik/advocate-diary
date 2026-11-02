package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import java.net.URISyntaxException;
import java.nio.file.Path;
import java.nio.file.Paths;
import org.testng.Assert;
import org.testng.annotations.Test;

public class FileUploadFlowTest extends BaseTest {
  private static final String FILE_NAME = "grounds-of-petition.pdf";

  @Test
  public void uploadWithoutFileShowsChooseAFile() {
    CasePage matter = openCasePage();
    Assert.assertTrue(matter.showsFileUpload());
    Assert.assertTrue(matter.acceptsPdfOnly());
    matter.clickUpload();
    Assert.assertTrue(matter.showsFileRequired());
  }

  @Test
  public void chooseFileAddLabelAndClickUploadShowsPaper() {
    CasePage matter = openCasePage();
    Assert.assertTrue(matter.showsFileUpload());
    matter.chooseFile(fixture());
    matter.selectLabel("petition");
    matter.clickUpload();
    matter.waitForPaper(FILE_NAME, "petition");
    Assert.assertTrue(matter.showsPaper(FILE_NAME, "petition"));
  }

  private Path fixture() {
    try {
      return Paths.get(getClass().getResource("/files/" + FILE_NAME).toURI());
    } catch (URISyntaxException e) {
      throw new IllegalStateException("missing upload fixture " + FILE_NAME, e);
    }
  }
}
