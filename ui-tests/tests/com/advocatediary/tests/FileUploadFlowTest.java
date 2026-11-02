package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import java.net.URISyntaxException;
import java.nio.file.Path;
import java.nio.file.Paths;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class FileUploadFlowTest extends BaseTest {
  private static final String FILE_NAME = "grounds-of-petition.txt";

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void uploadWithoutFileShowsChooseAFile() {
    CasePage matter = openCaseForUpload();
    Assert.assertTrue(matter.showsFileUpload(), "File, Label, and Upload should show under Papers");

    matter.clickUpload();
    Assert.assertTrue(matter.showsFileRequired());
  }

  @Test
  public void chooseFileAddLabelAndClickUploadShowsPaper() {
    CasePage matter = openCaseForUpload();
    Assert.assertTrue(matter.showsFileUpload(), "File, Label, and Upload should show under Papers");

    matter.chooseFile(fixture());
    matter.selectLabel("petition");
    matter.clickUpload();
    matter.waitForPaper(FILE_NAME, "petition");

    Assert.assertTrue(matter.showsPaper(FILE_NAME, "petition"));
  }

  private CasePage openCaseForUpload() {
    ClientPage client = signInToDiary()
        .nav()
        .openClients()
        .addClient()
        .saveClientWithNotes("Meera Rao", "9876543210", "meera@example.test", "Chennai", "Family matter");
    return client.openCase(
        "Rao v State",
        "Madras High Court",
        "civil",
        "WP 12 of 2026",
        "State of Tamil Nadu");
  }

  private Path fixture() {
    try {
      return Paths.get(getClass().getResource("/files/" + FILE_NAME).toURI());
    } catch (URISyntaxException e) {
      throw new IllegalStateException("missing upload fixture " + FILE_NAME, e);
    }
  }
}
